package pkg

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/otiai10/copy"
	monitoringV1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/postrender"

	imgErrors "github.com/d2iq-labs/helm-list-images/pkg/errors"
	"github.com/d2iq-labs/helm-list-images/pkg/k8s"
)

const (
	// ImageRegex is the default regex, that is used to split one big helm template to multiple templates.
	// Splitting templates eases the task of  identifying Kubernetes objects.
	ImageRegex = `---\n# Source:\s.*.`
)

// Images represents GetImages.
type Images struct {
	// Registries are list of registry names which we have filter out from
	Registries             []string
	Kind                   []string
	Values                 []string
	StringValues           []string
	FileValues             []string
	ExtraImagesFiles       []string
	ImageRegex             string
	ValueFiles             ValueFiles
	LogLevel               string
	FromRelease            bool
	UniqueImages           bool
	KubeVersion            string
	ChartVersionConstraint string
	PostRenderer           postrender.PostRenderer
	JSON                   bool
	YAML                   bool
	Table                  bool
	release                string
	chart                  string
	log                    *logrus.Logger
	writer                 *bufio.Writer
}

func (image *Images) SetRelease(release string) {
	image.release = release
}

func (image *Images) SetChart(chart string) {
	image.chart = chart
}

func (image *Images) SetWriter(writer io.Writer) {
	image.writer = bufio.NewWriter(writer)
}

func (image *Images) GetRelease() string {
	return image.release
}

func (image *Images) GetChart() string {
	return image.chart
}

func (image *Images) GetWriter() *bufio.Writer {
	return image.writer
}

// GetImages fetches all available images from the specified chart.
// Also filters identified images, to get just unique ones.
//
//nolint:gocyclo // Just a long function...
func (image *Images) GetImages() error {
	image.log.Debugf(
		"got all required values to fetch the images from chart/release '%s' proceeding further to fetch the same",
		image.release,
	)

	chart, err := image.getChartManifests()
	if err != nil {
		return err
	}

	image.log.Debugf("Rendered templates: %s", string(chart))

	images := make([]*k8s.Image, 0)
	kubeKindTemplates := image.GetTemplates(chart)

	for _, kubeKindTemplate := range kubeKindTemplates {
		currentKind, err := k8s.NewKind().Get(kubeKindTemplate)
		if err != nil {
			return err
		}

		if !funk.Contains(image.Kind, currentKind) {
			image.log.Debugf("either helm-list-images plugin does not support kind '%s' "+
				"at the moment or manifest might not have images to filter", currentKind,
			)

			continue
		}

		image.log.Debugf("fetching images from kind '%s'", currentKind)

		switch currentKind {
		case k8s.KindDeployment:
			deployImages, err := k8s.NewDeployment().Get(kubeKindTemplate)
			if err != nil {
				return err
			}

			images = append(images, deployImages)
		case k8s.KindStatefulSet:
			stsImages, err := k8s.NewStatefulSet().Get(kubeKindTemplate)
			if err != nil {
				return err
			}

			images = append(images, stsImages)
		case k8s.KindDaemonSet:
			daemonImages, err := k8s.NewDaemonSet().Get(kubeKindTemplate)
			if err != nil {
				return err
			}

			images = append(images, daemonImages)
		case k8s.KindReplicaSet:
			replicaSets, err := k8s.NewReplicaSets().Get(kubeKindTemplate)
			if err != nil {
				return err
			}

			images = append(images, replicaSets)
		case k8s.KindPod:
			pods, err := k8s.NewPod().Get(kubeKindTemplate)
			if err != nil {
				return err
			}

			images = append(images, pods)
		case k8s.KindCronJob:
			cronJob, err := k8s.NewCronjob().Get(kubeKindTemplate)
			if err != nil {
				return err
			}

			images = append(images, cronJob)
		case k8s.KindJob:
			job, err := k8s.NewJob().Get(kubeKindTemplate)
			if err != nil {
				return err
			}

			images = append(images, job)
		case monitoringV1.AlertmanagersKind:
			alertManager, err := k8s.NewAlertManager().Get(kubeKindTemplate)
			if err != nil {
				return err
			}

			images = append(images, alertManager)
		case monitoringV1.PrometheusesKind:
			prometheus, err := k8s.NewPrometheus().Get(kubeKindTemplate)
			if err != nil {
				return err
			}

			images = append(images, prometheus)
		case monitoringV1.ThanosRulerKind:
			thanosRuler, err := k8s.NewThanosRuler().Get(kubeKindTemplate)
			if err != nil {
				return err
			}

			images = append(images, thanosRuler)
		case k8s.KindThanos:
			thanos, err := k8s.NewThanos().Get(kubeKindTemplate)
			if err != nil {
				return err
			}

			images = append(images, thanos)
		case k8s.KindThanosReceiver:
			thanosReceiver, err := k8s.NewThanosReceiver().Get(kubeKindTemplate)
			if err != nil {
				return err
			}

			images = append(images, thanosReceiver)
		case k8s.KindGrafana:
			grafana, err := k8s.NewGrafana().Get(kubeKindTemplate)

			grafanaErr := &imgErrors.GrafanaAPIVersionSupportError{}
			if err != nil {
				if errors.As(err, &grafanaErr) {
					image.log.Debugf(
						"fetching images from Kind Grafana errored with %s",
						err.Error(),
					)

					continue
				} else {
					return err
				}
			}

			images = append(images, grafana)
		default:
			image.log.Debugf("kind '%s' is not supported at the moment", currentKind)
		}
	}

	return image.render(images)
}

func (image *Images) getChartManifests() ([]byte, error) {
	if image.FromRelease {
		image.log.Debugf(
			"from-release is selected, hence fetching manifests for '%s' from helm release",
			image.release,
		)

		return image.GetImagesFromRelease()
	}

	image.log.Debugf(
		"fetching manifests for '%s' by rendering helm template locally",
		image.release,
	)

	return image.getChartTemplate()
}

func (image *Images) getChartTemplate() ([]byte, error) {
	settings := cli.New()
	if strings.ToLower(image.LogLevel) == logrus.DebugLevel.String() {
		settings.Debug = true
	}

	actionConfig := new(action.Configuration)

	err := actionConfig.Init(
		settings.RESTClientGetter(),
		settings.Namespace(),
		os.Getenv("HELM_DRIVER"),
		image.log.Debugf,
	)
	if err != nil {
		image.log.Error("oops initialising helm client errored with", err)

		return nil, err
	}

	tmpDir, err := os.MkdirTemp("", ".helm-list-images-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	actualChartDir := tmpDir

	_, err = os.Stat(image.chart)

	switch {
	case err == nil:
		if err := copy.Copy(image.chart, actualChartDir); err != nil {
			return nil, fmt.Errorf("failed to copy chart to temporary directory: %w", err)
		}
	case filepath.IsAbs(image.chart):
		return nil, fmt.Errorf("specified chart path does not exist: %w", err)
	default:
		pull := action.NewPullWithOpts(action.WithConfig(actionConfig))
		pull.Settings = settings
		pull.Untar = true
		pull.DestDir = tmpDir
		pull.Version = image.ChartVersionConstraint

		out, err := pull.Run(image.chart)
		if err != nil {
			return nil, fmt.Errorf("failed to pull chart: %w (output: %s)", err, out)
		}

		actualChartDir = filepath.Join(tmpDir, filepath.Base(image.chart))
	}

	for fileIdx, extraImagesFile := range image.ExtraImagesFiles {
		extraImagesFileReader, err := os.Open(extraImagesFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open extra images file: %w", err)
		}

		resIdx := 0

		scanner := bufio.NewScanner(extraImagesFileReader)
		for scanner.Scan() {
			err = os.WriteFile(
				filepath.Join(
					actualChartDir,
					"templates",
					fmt.Sprintf("helm-list-images-extra-images-file-%d-%d.yaml", fileIdx, resIdx),
				),
				[]byte(`---
apiVersion: v1
kind: Pod
metadata:
  name: unused`+strconv.Itoa(fileIdx)+`
spec:
  containers:
  - image: `+scanner.Text()),
				0o400,
			)

			if err != nil {
				_ = extraImagesFileReader.Close()

				return nil, fmt.Errorf("failed to write extra images template: %w", err)
			}
			resIdx++
		}

		_ = extraImagesFileReader.Close()

		if scanner.Err() != nil {
			return nil, fmt.Errorf("failed to read extra images file: %w", err)
		}
	}

	chrt, err := loader.LoadDir(actualChartDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load chart: %w", err)
	}

	valueOpts := &values.Options{
		ValueFiles:   image.ValueFiles,
		StringValues: image.StringValues,
		Values:       image.Values,
		FileValues:   image.FileValues,
	}

	templateClient := action.NewInstall(actionConfig)
	templateClient.DryRun = true
	templateClient.ReleaseName = "release-name"
	templateClient.Replace = true // Skip the name check
	templateClient.ClientOnly = true
	templateClient.PostRenderer = image.PostRenderer

	if image.KubeVersion != "" {
		parsedKubeVersion, err := chartutil.ParseKubeVersion(image.KubeVersion)
		if err != nil {
			return nil, fmt.Errorf("invalid kube version '%s': %w", image.KubeVersion, err)
		}

		templateClient.KubeVersion = parsedKubeVersion
	}

	p := getter.All(settings)

	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return nil, fmt.Errorf("failed to merge values: %w", err)
	}

	templateOutput, err := templateClient.Run(chrt, vals)
	if err != nil {
		return nil, fmt.Errorf("failed to generate manifests: %w", err)
	}

	var manifests bytes.Buffer

	fmt.Fprintln(&manifests, strings.TrimSpace(templateOutput.Manifest))

	for _, h := range templateOutput.Hooks {
		if h != nil {
			fmt.Fprintf(&manifests, "---\n# Source: %s\n%s\n", h.Path, h.Manifest)
		}
	}

	return manifests.Bytes(), nil
}

func (image *Images) GetTemplates(template []byte) []string {
	image.log.Debugf("splitting helm manifests with regex pattern: '%s'", image.ImageRegex)
	temp := regexp.MustCompile(image.ImageRegex)
	kinds := temp.Split(string(template), -1)
	// Removing empty string at the beginning as splitting string always adds it in front.
	kinds = kinds[1:]

	return kinds
}

func GetImagesFromKind(kinds []*k8s.Image) []string {
	var images []string
	for _, knd := range kinds {
		images = append(images, knd.Image...)
	}

	return images
}

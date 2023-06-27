package k8s

import (
	"fmt"

	thanosalphav1 "github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/ghodss/yaml"
	grafanabetav1 "github.com/grafana-operator/grafana-operator/api/v1beta1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"

	imgErrors "github.com/d2iq-labs/helm-list-images/pkg/errors"
)

const (
	KindDeployment     = "Deployment"
	KindStatefulSet    = "StatefulSet"
	KindDaemonSet      = "DaemonSet"
	KindCronJob        = "CronJob"
	KindJob            = "Job"
	KindReplicaSet     = "ReplicaSet"
	KindPod            = "Pod"
	KindGrafana        = "Grafana"
	KindThanos         = "Thanos"
	KindThanosReceiver = "Receiver"
	kubeKind           = "kind"
)

type (
	Deployments  appsv1.Deployment
	StatefulSets appsv1.StatefulSet
	DaemonSets   appsv1.DaemonSet
	ReplicaSets  appsv1.ReplicaSet
	CronJob      batchv1.CronJob
	Job          batchv1.Job
	Pod          corev1.Pod
	Kind         map[string]interface{}
	containers   struct {
		containers []corev1.Container
	}
	AlertManager   monitoringv1.Alertmanager
	Prometheus     monitoringv1.Prometheus
	ThanosRuler    monitoringv1.ThanosRuler
	Grafana        grafanabetav1.Grafana
	Thanos         thanosalphav1.Thanos
	ThanosReceiver thanosalphav1.Receiver
)

type KindInterface interface {
	Get(dataMap string) (string, error)
}

type ImagesInterface interface {
	Get(dataMap string) (*Image, error)
}

type Image struct {
	Kind  string   `json:"kind,omitempty"`
	Name  string   `json:"name,omitempty"`
	Image []string `json:"image,omitempty"`
}

func (kin *Kind) Get(dataMap string) (string, error) {
	var kindYaml map[string]interface{}

	if err := yaml.Unmarshal([]byte(dataMap), &kindYaml); err != nil {
		return "", err
	}

	if len(kindYaml) != 0 {
		value, ok := kindYaml[kubeKind].(string)
		if !ok {
			return "", &imgErrors.ImageError{
				Message: "failed to get name from the manifest, 'kind' is not type string",
			}
		}

		return value, nil
	}

	return "", nil
}

func (dep *Deployments) Get(dataMap string) (*Image, error) {
	if err := yaml.Unmarshal([]byte(dataMap), &dep); err != nil {
		return nil, err
	}

	depContainers := containers{
		append(dep.Spec.Template.Spec.Containers, dep.Spec.Template.Spec.InitContainers...),
	}

	images := &Image{
		Kind:  KindDeployment,
		Name:  dep.Name,
		Image: depContainers.getImages(),
	}

	return images, nil
}

func (dep *StatefulSets) Get(dataMap string) (*Image, error) {
	if err := yaml.Unmarshal([]byte(dataMap), &dep); err != nil {
		return nil, err
	}

	depContainers := containers{
		append(dep.Spec.Template.Spec.Containers, dep.Spec.Template.Spec.InitContainers...),
	}

	images := &Image{
		Kind:  KindStatefulSet,
		Name:  dep.Name,
		Image: depContainers.getImages(),
	}

	return images, nil
}

func (dep *DaemonSets) Get(dataMap string) (*Image, error) {
	if err := yaml.Unmarshal([]byte(dataMap), &dep); err != nil {
		return nil, err
	}

	depContainers := containers{
		append(dep.Spec.Template.Spec.Containers, dep.Spec.Template.Spec.InitContainers...),
	}

	images := &Image{
		Kind:  KindDaemonSet,
		Name:  dep.Name,
		Image: depContainers.getImages(),
	}

	return images, nil
}

func (dep *CronJob) Get(dataMap string) (*Image, error) {
	if err := yaml.Unmarshal([]byte(dataMap), &dep); err != nil {
		return nil, err
	}

	depContainers := containers{append(dep.Spec.JobTemplate.Spec.Template.Spec.Containers,
		dep.Spec.JobTemplate.Spec.Template.Spec.InitContainers...)}

	images := &Image{
		Kind:  KindCronJob,
		Name:  dep.Name,
		Image: depContainers.getImages(),
	}

	return images, nil
}

func (dep *Job) Get(dataMap string) (*Image, error) {
	if err := yaml.Unmarshal([]byte(dataMap), &dep); err != nil {
		return nil, err
	}

	depContainers := containers{
		append(dep.Spec.Template.Spec.Containers, dep.Spec.Template.Spec.InitContainers...),
	}

	images := &Image{
		Kind:  KindJob,
		Name:  dep.Name,
		Image: depContainers.getImages(),
	}

	return images, nil
}

func (dep *ReplicaSets) Get(dataMap string) (*Image, error) {
	if err := yaml.Unmarshal([]byte(dataMap), &dep); err != nil {
		return nil, err
	}

	depContainers := containers{
		append(dep.Spec.Template.Spec.Containers, dep.Spec.Template.Spec.InitContainers...),
	}

	images := &Image{
		Kind:  KindReplicaSet,
		Name:  dep.Name,
		Image: depContainers.getImages(),
	}

	return images, nil
}

func (dep *Pod) Get(dataMap string) (*Image, error) {
	if err := yaml.Unmarshal([]byte(dataMap), &dep); err != nil {
		return nil, err
	}

	depContainers := containers{append(dep.Spec.Containers, dep.Spec.InitContainers...)}

	images := &Image{
		Kind:  KindPod,
		Name:  dep.Name,
		Image: depContainers.getImages(),
	}

	return images, nil
}

func (dep *AlertManager) Get(dataMap string) (*Image, error) {
	if err := yaml.Unmarshal([]byte(dataMap), &dep); err != nil {
		return nil, err
	}

	images := &Image{
		Kind:  monitoringv1.AlertmanagersKind,
		Name:  dep.Name,
		Image: []string{*dep.Spec.Image},
	}

	return images, nil
}

func (dep *Prometheus) Get(dataMap string) (*Image, error) {
	if err := yaml.Unmarshal([]byte(dataMap), &dep); err != nil {
		return nil, err
	}

	var imageNames []string

	depContainers := containers{append(dep.Spec.Containers, dep.Spec.InitContainers...)}

	imageNames = append(imageNames, depContainers.getImages()...)
	imageNames = append(imageNames, *dep.Spec.Image)

	images := &Image{
		Kind:  monitoringv1.PrometheusesKind,
		Name:  dep.Name,
		Image: imageNames,
	}

	return images, nil
}

func (dep *ThanosRuler) Get(dataMap string) (*Image, error) {
	if err := yaml.Unmarshal([]byte(dataMap), &dep); err != nil {
		return nil, err
	}

	var imageNames []string

	depContainers := containers{append(dep.Spec.Containers, dep.Spec.InitContainers...)}

	imageNames = append(imageNames, depContainers.getImages()...)
	imageNames = append(imageNames, dep.Spec.Image)

	images := &Image{
		Kind:  monitoringv1.ThanosRulerKind,
		Name:  dep.Name,
		Image: imageNames,
	}

	return images, nil
}

func (dep *Grafana) Get(dataMap string) (*Image, error) {
	if err := yaml.Unmarshal([]byte(dataMap), &dep); err != nil {
		return nil, err
	}

	if dep.APIVersion == "integreatly.org/v1alpha1" {
		return nil, &imgErrors.GrafanaAPIVersionSupportError{
			Message: fmt.Sprintf(
				"plugin supports the latest api version and '%s' is not supported",
				dep.APIVersion,
			),
		}
	}

	grafanaDeployment := dep.Spec.Deployment.Spec.Template.Spec
	depContainers := containers{
		append(grafanaDeployment.Containers, grafanaDeployment.InitContainers...),
	}

	images := &Image{
		Kind:  KindGrafana,
		Name:  dep.Name,
		Image: depContainers.getImages(),
	}

	return images, nil
}

func (dep *Thanos) Get(dataMap string) (*Image, error) {
	if err := yaml.Unmarshal([]byte(dataMap), &dep); err != nil {
		return nil, err
	}

	var thanosContainers []corev1.Container
	thanosContainers = append(thanosContainers,
		dep.Spec.Rule.StatefulsetOverrides.Spec.Template.Spec.Containers...)
	thanosContainers = append(thanosContainers,
		dep.Spec.Rule.StatefulsetOverrides.Spec.Template.Spec.InitContainers...)
	thanosContainers = append(thanosContainers,
		dep.Spec.Query.DeploymentOverrides.Spec.Template.Spec.Containers...)
	thanosContainers = append(thanosContainers,
		dep.Spec.Query.DeploymentOverrides.Spec.Template.Spec.InitContainers...)
	thanosContainers = append(thanosContainers,
		dep.Spec.StoreGateway.DeploymentOverrides.Spec.Template.Spec.Containers...)
	thanosContainers = append(thanosContainers,
		dep.Spec.StoreGateway.DeploymentOverrides.Spec.Template.Spec.InitContainers...)
	thanosContainers = append(thanosContainers,
		dep.Spec.QueryFrontend.DeploymentOverrides.Spec.Template.Spec.Containers...)
	thanosContainers = append(thanosContainers,
		dep.Spec.QueryFrontend.DeploymentOverrides.Spec.Template.Spec.InitContainers...)

	depContainers := containers{thanosContainers}

	images := &Image{
		Kind:  KindThanos,
		Name:  dep.Name,
		Image: depContainers.getImages(),
	}

	return images, nil
}

func (dep *ThanosReceiver) Get(dataMap string) (*Image, error) {
	if err := yaml.Unmarshal([]byte(dataMap), &dep); err != nil {
		return nil, err
	}

	var receiverGroupContainers []corev1.Container
	for idx := range dep.Spec.ReceiverGroups {
		receiverGroup := dep.Spec.ReceiverGroups[idx]
		receiverGroupContainers = append(receiverGroupContainers,
			receiverGroup.StatefulSetOverrides.Spec.Template.Spec.Containers...)
		receiverGroupContainers = append(receiverGroupContainers,
			receiverGroup.StatefulSetOverrides.Spec.Template.Spec.InitContainers...)
	}

	depContainers := containers{receiverGroupContainers}

	images := &Image{
		Kind:  KindThanosReceiver,
		Name:  dep.Name,
		Image: depContainers.getImages(),
	}

	return images, nil
}

func NewDeployment() ImagesInterface {
	return &Deployments{}
}

func NewStatefulSet() ImagesInterface {
	return &StatefulSets{}
}

func NewDaemonSet() ImagesInterface {
	return &DaemonSets{}
}

func NewReplicaSets() ImagesInterface {
	return &ReplicaSets{}
}

func NewCronjob() ImagesInterface {
	return &CronJob{}
}

func NewJob() ImagesInterface {
	return &Job{}
}

func NewPod() ImagesInterface {
	return &Pod{}
}

func NewAlertManager() ImagesInterface {
	return &AlertManager{}
}

func NewPrometheus() ImagesInterface {
	return &Prometheus{}
}

func NewThanosRuler() ImagesInterface {
	return &ThanosRuler{}
}

func NewGrafana() ImagesInterface {
	return &Grafana{}
}

func NewThanos() ImagesInterface {
	return &Thanos{}
}

func NewThanosReceiver() ImagesInterface {
	return &ThanosReceiver{}
}

func NewKind() KindInterface {
	return &Kind{}
}

func SupportedKinds() []string {
	kinds := []string{
		KindDeployment, KindStatefulSet, KindDaemonSet,
		KindCronJob, KindJob, KindReplicaSet, KindPod,
		monitoringv1.AlertmanagersKind, monitoringv1.PrometheusesKind, monitoringv1.ThanosRulerKind,
		KindGrafana, KindThanos, KindThanosReceiver,
	}
	return kinds
}

func (cont containers) getImages() []string {
	images := make([]string, 0, len(cont.containers))
	for idx := range cont.containers {
		images = append(images, cont.containers[idx].Image)
	}

	return images
}

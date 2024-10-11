package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/component-base/version/verflag"

	"github.com/d2iq-labs/helm-list-images/pkg"
	"github.com/d2iq-labs/helm-list-images/pkg/k8s"
)

// Registers all global flags to utility itself.
func registerFlags(cmd *cobra.Command, images *pkg.Images) {
	verflag.AddFlags(cmd.Flags())

	cmd.Flags().StringArrayVar(&images.Values, "set", []string{},
		"set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	cmd.Flags().StringArrayVar(&images.StringValues, "set-string", []string{},
		"set STRING values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	cmd.Flags().StringArrayVar(&images.FileValues, "set-file", []string{},
		"set values from respective files specified via the command line (can specify multiple or separate values with "+
			"commas: key1=path1,key2=path2)")
	cmd.Flags().VarP(&images.ValueFiles, "values", "f",
		"specify values in a YAML file (can specify multiple)")
	cmd.Flags().
		StringArrayVar(
			&images.JSONValues,
			"set-json",
			[]string{},
			`set JSON values on the command line (can specify multiple or separate values with commas: `+
				`key1=jsonval1,key2=jsonval2)`,
		)
	cmd.Flags().
		StringArrayVar(&images.LiteralValues, "set-literal", []string{}, "set a literal STRING value on the command line")

	cmd.Flags().StringSliceVarP(&images.Registries, "registry", "r", nil,
		"registry name (docker images belonging to this registry)")
	cmd.Flags().StringSliceVarP(&images.Kind, "kind", "k", k8s.SupportedKinds(),
		"kubernetes app kind to fetch the images from")
	cmd.Flags().StringVarP(&images.LogLevel, "log-level", "l", "info",
		"log level for the plugin helm list-images (defaults to info)")
	cmd.Flags().StringVarP(&images.ImageRegex, "image-regex", "", pkg.ImageRegex,
		"regex used to split helm template rendered")
	cmd.Flags().BoolVarP(&images.UniqueImages, "unique", "u", true,
		"enable the flag if duplicates to be removed from the retrieved list")
	cmd.Flags().BoolVarP(&images.SortImages, "sort", "s", true,
		"enable the flag to sort images")
	cmd.Flags().BoolVarP(&images.JSON, "json", "j", false,
		"enable the flag to display images retrieved in json format (disabled by default)")
	cmd.Flags().BoolVarP(&images.YAML, "yaml", "y", false,
		"enable the flag to display images retrieved in yaml format (disabled by default)")
	cmd.Flags().BoolVarP(&images.Table, "table", "t", false,
		"enable the flag to display images retrieved in table format (disabled by default)")
	cmd.Flags().BoolVarP(&images.FromRelease, "from-release", "", false,
		"enable the flag to fetch the images from release instead (disabled by default)")
	cmd.Flags().StringArrayVar(&images.ExtraImagesFiles, "extra-images-file", []string{},
		"optional Helm template files to derive extra images from")
	cmd.MarkFlagsMutuallyExclusive("from-release", "extra-images-file")
	cmd.Flags().StringVar(&images.KubeVersion, "kube-version", "1.29.0",
		"Kubernetes version used for Capabilities.KubeVersion when rendering charts")
	cmd.Flags().StringArrayVar(&images.APIVersions, "api-versions", []string{},
		"Kubernetes api versions used for Capabilities.APIVersions when rendering charts")
	cmd.Flags().StringVar(&images.ChartVersionConstraint, "chart-version", "",
		"specify a version constraint for the chart version to use. This constraint can be a specific tag (e.g. 1.1.1) or "+
			"it may reference a valid range (e.g. ^2.0.0). If this is not specified, the latest version is used",
	)
	cmd.Flags().StringVar(&images.RepoURL, "repo", "",
		"chart repository url where to locate the requested chart",
	)
	cmd.MarkFlagsMutuallyExclusive("from-release", "chart-version")
	cmd.MarkFlagsMutuallyExclusive("from-release", "repo")
	cmd.Flags().
		BoolVar(&images.IncludeTestImages, "include-test-images", false, "include images required for Helm tests")
}

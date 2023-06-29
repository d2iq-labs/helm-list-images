package cmd

import (
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/component-base/version/verflag"

	"github.com/d2iq-labs/helm-list-images/pkg"
	imgErrors "github.com/d2iq-labs/helm-list-images/pkg/errors"
)

func GetRootCommand() *cobra.Command {
	images := pkg.Images{}

	rootCommand := &cobra.Command{
		Use:   "list-images CHART|RELEASE [flags]",
		Short: "Fetches all images those are part of specified chart/release",
		Long: "Lists all images those are part of specified chart/release and matches the pattern or part of specified " +
			"registry.",
		Example: `  helm list-images path/to/chart/prometheus-standalone -f ~/path/to/override-config.yaml
  helm list-images prometheus-standalone --from-release --registry quay.io
  helm list-images prometheus-standalone --from-release --registry quay.io --unique
  helm list-images prometheus-standalone --from-release --registry quay.io --yaml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			verflag.PrintAndExitIfRequested()

			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err
			}

			log.SetOutput(io.Discard)
			images.SetLogger(images.LogLevel)
			images.SetWriter(os.Stdout)
			cmd.SilenceUsage = true

			if images.FromRelease {
				images.SetRelease(args[0])
			} else {
				images.SetChart(args[0])
			}

			if (images.JSON && images.YAML && images.Table) || (images.JSON && images.YAML) ||
				(images.Table && images.YAML) || (images.Table && images.JSON) {
				return &imgErrors.MultipleFormatError{
					Message: "cannot render the output to multiple format, enable any of '--yaml --json --table' at a time",
				}
			}

			return images.GetImages()
		},
	}

	registerFlags(rootCommand, &images)

	rootCommand.SetUsageTemplate(getUsageTemplate())
	rootCommand.DisableAutoGenTag = true

	return rootCommand
}

func getUsageTemplate() string {
	return `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if gt (len .Aliases) 0}}{{printf "\n" }}
Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}{{printf "\n" }}
Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{printf "\n"}}
Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}{{printf "\n"}}
Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}{{printf "\n"}}
Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}{{printf "\n"}}
Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}
{{if .HasAvailableSubCommands}}{{printf "\n"}}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
{{printf "\n"}}`
}

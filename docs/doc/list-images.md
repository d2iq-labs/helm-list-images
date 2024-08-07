## list-images

Fetches all images those are part of specified chart/release

### Synopsis

Lists all images those are part of specified chart/release and matches the pattern or part of specified registry.

```
list-images CHART|RELEASE [flags]
```

### Examples

```
  helm list-images path/to/chart/prometheus-standalone -f ~/path/to/override-config.yaml
  helm list-images prometheus --repository-url https://prometheus-community.github.io/helm-charts
  helm list-images prometheus-standalone --from-release --registry quay.io
  helm list-images prometheus-standalone --from-release --registry quay.io --unique
  helm list-images prometheus-standalone --from-release --registry quay.io --yaml
```

### Options

```
      --chart-version string            specify a version constraint for the chart version to use. This constraint can be a specific tag (e.g. 1.1.1) or it may reference a valid range (e.g. ^2.0.0). If this is not specified, the latest version is used
      --extra-images-file stringArray   optional Helm template files to derive extra images from
      --from-release                    enable the flag to fetch the images from release instead (disabled by default)
  -h, --help                            help for list-images
      --image-regex string              regex used to split helm template rendered (default "---\\n# Source:\\s.*.")
      --include-test-images             include images required for Helm tests
  -j, --json                            enable the flag to display images retrieved in json format (disabled by default)
  -k, --kind strings                    kubernetes app kind to fetch the images from (default [Deployment,StatefulSet,DaemonSet,CronJob,Job,ReplicaSet,Pod,Alertmanager,Prometheus,ThanosRuler,Grafana,Thanos,Receiver])
      --kube-version string             Kubernetes version used for Capabilities.KubeVersion when rendering charts (default "1.29.0")
  -l, --log-level string                log level for the plugin helm list-images (defaults to info) (default "info")
  -r, --registry strings                registry name (docker images belonging to this registry)
      --repo string                     chart repository url where to locate the requested chart
      --set stringArray                 set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)
      --set-file stringArray            set values from respective files specified via the command line (can specify multiple or separate values with commas: key1=path1,key2=path2)
      --set-json stringArray            set JSON values on the command line (can specify multiple or separate values with commas: key1=jsonval1,key2=jsonval2)
      --set-literal stringArray         set a literal STRING value on the command line
      --set-string stringArray          set STRING values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)
  -s, --sort                            enable the flag to sort images (default true)
  -t, --table                           enable the flag to display images retrieved in table format (disabled by default)
  -u, --unique                          enable the flag if duplicates to be removed from the retrieved list (default true)
  -f, --values ValueFiles               specify values in a YAML file (can specify multiple) (default [])
      --version version[=true]          --version, --version=raw prints version information and quits; --version=vX.Y.Z... sets the reported version
  -y, --yaml                            enable the flag to display images retrieved in yaml format (disabled by default)
```


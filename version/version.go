// Package version powers the versioning of terragen.
package version

var (
	// Version specifies the version of the application and cannot be changed by end user.
	Version string

	// BuildDate of the app.
	BuildDate string
	// Platform is the combination of OS and Architecture for which the binary is built for.
	Platform string
	// Revision represents the git revision used to build the current version of app.
	Revision string
)

// BuildInfo represents version of utility.
type BuildInfo struct {
	Version     string
	Revision    string
	Environment string
	BuildDate   string
	Platform    string
}

// GetBuildInfo return the version and other build info of the application.
func GetBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   Version,
		Revision:  Revision,
		Platform:  Platform,
		BuildDate: BuildDate,
	}
}

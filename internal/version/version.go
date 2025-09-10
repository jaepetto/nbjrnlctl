// Package version provides version information for the application.
// Date: 2025-09-10
// Description: Version tracking package that provides application version information
// and build metadata through compile-time injection.

package version

import (
	"fmt"
	"runtime"
)

// These variables are set at build time using ldflags.
// They provide version information and build metadata.
var (
	// Version is the semantic version of the application (e.g., "1.0.0")
	Version = "dev"

	// GitCommit is the git commit hash at build time
	GitCommit = "unknown"

	// BuildDate is the timestamp when the binary was built
	BuildDate = "unknown"

	// GoVersion is the Go version used to compile the binary
	GoVersion = runtime.Version()

	// Platform is the OS/architecture combination (e.g., "darwin/amd64")
	Platform = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)

// GetVersion returns the full version information as a formatted string.
func GetVersion() string {
	return fmt.Sprintf("nbjrnlctl version %s (commit: %s, built: %s, go: %s, platform: %s)",
		Version, GitCommit, BuildDate, GoVersion, Platform)
}

// GetShortVersion returns just the semantic version.
func GetShortVersion() string {
	return Version
}

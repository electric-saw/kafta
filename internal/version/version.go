package version

import (
	"runtime"
)

var (
	version = "v0.1.2"

	metadata  = ""
	gitCommit = ""
)

type BuildInfo struct {
	// Version is the current semver.
	Version string `json:"version,omitempty"`
	// GitCommit is the git sha1.
	GitCommit string `json:"git_commit,omitempty"`
	// GoVersion is the version of the Go compiler used.
	GoVersion string `json:"go_version,omitempty"`
}

func GetVersion() string {
	if metadata == "" {
		return version
	}
	return version + "+" + metadata
}

// Get returns build info
func Get() BuildInfo {
	v := BuildInfo{
		Version:   GetVersion(),
		GitCommit: gitCommit,
		GoVersion: runtime.Version(),
	}

	return v
}

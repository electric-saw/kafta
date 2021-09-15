package version

import (
	"flag"
	"runtime"
)

var (
	version = "v0.0.1"

	metadata     = ""
	gitCommit    = ""
	gitTreeState = ""
)

type BuildInfo struct {
	// Version is the current semver.
	Version string `json:"version,omitempty"`
	// GitCommit is the git sha1.
	GitCommit string `json:"git_commit,omitempty"`
	// GitTreeState is the state of the git tree.
	GitTreeState string `json:"git_tree_state,omitempty"`
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
		Version:      GetVersion(),
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		GoVersion:    runtime.Version(),
	}

	if flag.Lookup("test.v") != nil {
		v.GoVersion = ""
	}
	return v
}

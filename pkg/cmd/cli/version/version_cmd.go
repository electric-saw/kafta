package version

import (
	"fmt"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/version"
	"github.com/spf13/cobra"
)

func NewCmdVersion(config *configuration.Configuration) *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   fmt.Sprintf("Print the %s version", config.AppName),
		Aliases: []string{"v"},
		Run: func(cmd *cobra.Command, args []string) {
			printVersion()
		},
	}
}

func printVersion() {
	v := version.Get()
	fmt.Println("Version        ", v.Version)
	fmt.Println("Git commit     ", v.GitCommit)
	fmt.Println("Go version     ", v.GoVersion)
}

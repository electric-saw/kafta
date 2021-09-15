package cluster

import (
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/spf13/cobra"
)

func NewCmdCluster(config *configuration.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "cluster management",
	}

	cmd.AddCommand(NewCmdDescribeCluster(config))

	return cmd
}

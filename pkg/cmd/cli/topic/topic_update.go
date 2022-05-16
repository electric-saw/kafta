package topic

import (
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
)

type updateTopicOptions struct {
	config *configuration.Configuration
	name   string
	key    string
	value  string
}

func NewCmdUpdateTopic(config *configuration.Configuration) *cobra.Command {
	options := &updateTopicOptions{config: config}
	cmd := &cobra.Command{
		Use:   "update NAME KEY VALUE",
		Short: "update topics",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd))
			cmdutil.CheckErr(options.run())
		},
	}

	return cmd

}

func (o *updateTopicOptions) complete(cmd *cobra.Command) error {
	args := cmd.Flags().Args()
	if len(args) > 3 {
		return cmdutil.HelpErrorf(cmd, "Unexpected args: %v", args)
	}
	if len(args) == 3 {
		o.name = args[0]
		o.key = args[1]
		o.value = args[2]
	}
	return nil
}

func (o *updateTopicOptions) run() error {
	conn := kafka.MakeConnection(o.config)
	defer conn.Close()

	return kafka.SetProp(conn, o.name, o.key, o.value)
}

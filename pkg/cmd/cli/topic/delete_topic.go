package topic

import (
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/Songmu/prompter"
	"github.com/spf13/cobra"
)

type deleteTopicOptions struct {
	config *configuration.Configuration
	name   string
	quiet  bool
}

func NewCmdDeleteTopic(config *configuration.Configuration) *cobra.Command {
	options := &deleteTopicOptions{config: config}
	cmd := &cobra.Command{
		Use:   "delete NAME [--quiet]",
		Short: "Delete topics",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return ValidTopics(config, len(args) > 0)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd))
			cmdutil.CheckErr(options.run())
		},
	}

	cmd.Flags().BoolVarP(&options.quiet, "quiet", "q", false, "Quiet mode")

	return cmd

}

func (o *deleteTopicOptions) complete(cmd *cobra.Command) error {
	args := cmd.Flags().Args()
	if len(args) > 1 {
		return cmdutil.HelpErrorf(cmd, "Unexpected args: %v", args)
	}
	if len(args) == 1 {
		o.name = args[0]
	}
	return nil
}

func (o *deleteTopicOptions) run() error {
	if !o.quiet {
		if !prompter.YN("really want to delete?", false) {
			return nil
		}
	}

	conn := kafka.MakeConnection(o.config)
	defer conn.Close()
	return kafka.DeleteTopic(conn, o.name)
}

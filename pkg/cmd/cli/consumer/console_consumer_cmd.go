package consumer

import (
	"strings"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
)

type consumerOptions struct {
	config  *configuration.Configuration
	topic   string
	group   string
	verbose bool
}

func NewCmdConsumeMessage(config *configuration.Configuration) *cobra.Command {
	options := &consumerOptions{config: config}
	cmd := &cobra.Command{
		Use:   "consumer TOPIC [group=cgName] [--verbose]",
		Short: "Consume messages",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return ValidTopics(config, len(args) > 0)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd))
			cmdutil.CheckErr(options.run())
		},
	}
	cmd.Flags().BoolVarP(&options.verbose, "verbose", "v", false, "Verboset mode")
	return cmd
}

func (o *consumerOptions) complete(cmd *cobra.Command) error {
	args := cmd.Flags().Args()
	if len(args) > 2 {
		return cmdutil.HelpErrorf(cmd, "Unexpected args: %v", args)
	}

	if len(args) > 1 {
		o.topic = args[0]
		o.group = strings.Split(args[1], "=")[1]
	} else if len(args) == 1 {
		o.topic = args[0]
		o.group = "kafta-cli"
	}
	return nil
}

func (o *consumerOptions) run() error {
	conn := kafka.EstablishKafkaConnection(o.config)
	defer conn.Close()
	return kafka.ConsumeMessage(conn, o.topic, o.group, o.verbose)
}

func ValidTopics(
	config *configuration.Configuration,
	hasArgs bool,
) ([]string, cobra.ShellCompDirective) {
	var topicsList []string

	if hasArgs {
		return topicsList, cobra.ShellCompDirectiveNoFileComp
	}

	conn := kafka.EstablishKafkaConnection(config)
	topics := kafka.ListAllTopics(conn)
	for name := range topics {
		topicsList = append(topicsList, name)
	}

	return topicsList, cobra.ShellCompDirectiveNoFileComp
}

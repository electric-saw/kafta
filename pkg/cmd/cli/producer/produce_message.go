package producer

import (
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
)

type producerOptions struct {
	config *configuration.Configuration
	topic  string
}

func NewCmdProduceMessage(config *configuration.Configuration) *cobra.Command {
	options := &producerOptions{config: config}
	cmd := &cobra.Command{
		Use:   "producer TOPIC",
		Short: "Produce messages",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return ValidTopics(config, len(args) > 0)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd))
			cmdutil.CheckErr(options.run())
		},
	}

	return cmd
}

func (o *producerOptions) complete(cmd *cobra.Command) error {
	args := cmd.Flags().Args()
	if len(args) > 2 {
		return cmdutil.HelpErrorf(cmd, "Unexpected args: %v", args)
	}
	if len(args) == 1 {
		o.topic = args[0]
	}
	return nil
}

func (o *producerOptions) run() error {
	conn := kafka.EstablishKafkaConnection(o.config)
	defer conn.Close()
	return kafka.ProduceMessage(conn, o.topic)
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

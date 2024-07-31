package topic

import (
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	"github.com/spf13/cobra"
)

func NewCmdTopic(config *configuration.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "topic",
		Short: "Topics management",
	}

	cmd.AddCommand(NewCmdCreateTopic(config))
	cmd.AddCommand(NewCmdDeleteTopic(config))
	cmd.AddCommand(NewCmdListTopic(config))
	cmd.AddCommand(NewCmdDescribeTopic(config))
	cmd.AddCommand(NewCmdConfigUpdateTopic(config))
	cmd.AddCommand(NewCmdConfigResetTopic(config))
	cmd.AddCommand(NewCmdListConfigs(config))

	return cmd
}

func ValidTopics(config *configuration.Configuration, hasArgs bool) ([]string, cobra.ShellCompDirective) {
	var topicsList []string

	if hasArgs {
		return topicsList, cobra.ShellCompDirectiveNoFileComp
	}

	conn := kafka.MakeConnection(config)
	topics := kafka.ListAllTopics(conn)
	for name := range topics {
		topicsList = append(topicsList, name)
	}

	return topicsList, cobra.ShellCompDirectiveNoFileComp
}

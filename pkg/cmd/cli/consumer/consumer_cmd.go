package consumer

import (
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	"github.com/spf13/cobra"
)

// Lag
// offset manage

func NewCmdConsumer(config *configuration.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "consumer",
		Short: "Consumer group management",
	}

	cmd.AddCommand(NewCmdListConsumer(config))
	cmd.AddCommand(NewCmdDescribeConsumer(config))
	cmd.AddCommand(NewCmdLagConsumer(config))
	cmd.AddCommand(NewCmdDeleteConsumer(config))
	cmd.AddCommand(NewCmdResetOffset(config))

	return cmd
}

func ValidConsumers(
	config *configuration.Configuration,
	hasArgs bool,
) ([]string, cobra.ShellCompDirective) {
	var consumersList []string
	conn := kafka.EstablishKafkaConnection(config)
	consumers := kafka.ListConsumerGroups(conn)
	for name := range consumers {
		consumersList = append(consumersList, name)
	}

	return consumersList, cobra.ShellCompDirectiveNoFileComp
}

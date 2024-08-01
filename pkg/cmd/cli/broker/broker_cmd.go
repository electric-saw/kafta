package broker

import (
	"strconv"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	"github.com/spf13/cobra"
)

func NewCmdBroker(config *configuration.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "broker",
		Short: "broker management",
	}

	// Get config
	// Get all configs
	// Update config
	// Reset config?
	cmd.AddCommand(NewCmdClusterGetConfig(config))
	cmd.AddCommand(NewCmdClusterDescribe(config))

	return cmd
}

func ValidBrokers(config *configuration.Configuration, hasArgs bool) ([]string, cobra.ShellCompDirective) {
	var brokerList []string
	conn := kafka.MakeConnection(config)
	brokers := kafka.GetBrokers(conn)
	for _, broker := range brokers {
		brokerList = append(brokerList, broker.Host, strconv.Itoa(int(broker.Id)))
	}

	return brokerList, cobra.ShellCompDirectiveNoFileComp
}

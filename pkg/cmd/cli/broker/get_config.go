package broker

import (
	"strconv"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	"github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

type clusterConfig struct {
	config   *configuration.Configuration
	brokerId string
}

func NewCmdClusterGetConfig(config *configuration.Configuration) *cobra.Command {
	options := &clusterConfig{config: config}
	cmd := &cobra.Command{
		Use:   "get-configs BROKER_ID (not required)",
		Short: "Show broker configs, by default is used coodinator",
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(options.defaultValue(cmd))
			options.run(cmd)
		},
	}

	return cmd
}

func (o *clusterConfig) defaultValue(cmd *cobra.Command) error {
	args := cmd.Flags().Args()

	if len(args) == 0 {
		conn := kafka.EstablishKafkaConnection(o.config)
		defer conn.Close()

		brokers := kafka.GetBrokers(conn)

		for _, broker := range brokers {
			if broker.IsController {
				o.brokerId = strconv.FormatInt(int64(broker.ID()), 10)
				break
			}
		}

		if o.brokerId == "" {
			return util.HelpError(cmd, "Impossible find BrokerId coordinator")
		}
	} else {
		o.brokerId = args[0]
	}

	return nil
}

func (o *clusterConfig) run(_ *cobra.Command) {
	conn := kafka.EstablishKafkaConnection(o.config)
	defer conn.Close()

	configs := kafka.DescribeBrokerConfig(conn, o.brokerId)

	header := table.Row{"name", "value", "default"}
	rows := []table.Row{}

	for _, config := range configs {
		rows = append(
			rows,
			table.Row{config.Name, text.WrapText(config.Value, 100), config.Default},
		)
	}

	util.PrintTable(header, rows)
}

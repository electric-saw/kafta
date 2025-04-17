package consumer

import (
	"sort"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	"github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type listConsumerOptions struct {
	config *configuration.Configuration
}

func NewCmdListConsumer(config *configuration.Configuration) *cobra.Command {
	options := &listConsumerOptions{config: config}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List consumers",
		Run: func(cmd *cobra.Command, args []string) {
			options.run()
		},
	}

	return cmd
}

func (o *listConsumerOptions) run() {
	conn := kafka.EstablishKafkaConnection(o.config)
	defer conn.Close()
	consumers := kafka.ListConsumerGroupDescriptions(conn)
	rows := []table.Row{}

	keys := make([]string, 0, len(consumers))
	for k := range consumers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		consumerType := consumers[name]
		rows = append(rows, table.Row{name, consumerType.ProtocolType, consumerType.State})
	}

	util.PrintTable(table.Row{"name", "type", "state"}, rows)
}

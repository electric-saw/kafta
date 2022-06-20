package cluster

import (
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	"github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type describeClusters struct {
	config *configuration.Configuration
}

func NewCmdDescribeCluster(config *configuration.Configuration) *cobra.Command {
	options := &describeClusters{config: config}
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe current-cluster",
		Run: func(cmd *cobra.Command, args []string) {
			options.run()
		},
	}
	return cmd
}

func (o *describeClusters) run() {
	conn := kafka.MakeConnection(o.config)
	defer conn.Close()
	brokers := kafka.GetBrokers(conn)

	header := table.Row{"id", "address", "rack", "controller"}
	rows := []table.Row{}

	for _, broker := range brokers {
		rows = append(rows, table.Row{broker.ID(), broker.Address, broker.Broker.Rack(), broker.IsController})
	}

	util.PrintTable(header, rows)
}

package broker

import (
	"strconv"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	"github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type clusterDescribe struct {
	config   *configuration.Configuration
	brokerId string
}

func NewCmdClusterDescribe(config *configuration.Configuration) *cobra.Command {
	options := &clusterDescribe{config: config}
	cmd := &cobra.Command{
		Use:   "describe BROKER_ID (not required)",
		Short: "Show broker details",
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(options.defaultValue(cmd))
			options.run(cmd)
		},
	}

	return cmd
}

func (o *clusterDescribe) defaultValue(cmd *cobra.Command) error {
	args := cmd.Flags().Args()

	if len(args) == 0 {
		conn := kafka.MakeConnection(o.config)
		defer conn.Close()

		brokers := kafka.GetBrokers(conn)

		for _, broker := range brokers {
			if broker.IsController {
				o.brokerId = strconv.Itoa(int(broker.ID()))
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

func (o *clusterDescribe) run(cmd *cobra.Command) {
	conn := kafka.MakeConnection(o.config)
	defer conn.Close()

	description := kafka.DescribeBroker(conn, o.brokerId)

	// Print broker details
	header := table.Row{"", ""}
	rows := []table.Row{}

	rows = append(rows, table.Row{"ID", description.Details.ID()})
	rows = append(rows, table.Row{"Host", description.Details.Host})
	rows = append(rows, table.Row{"Rack", description.Details.Rack})
	rows = append(rows, table.Row{"Controller", description.Details.IsController})

	util.PrintTable(header, rows)

	// Print broker logs
	header = table.Row{"Path", "Permanent", "Temporary", "Total"}
	rows = []table.Row{}

	for _, log := range description.Logs {
		for _, entry := range log.Entries {
			rows = append(rows, table.Row{log.Path, entry.Permanent, entry.Temporary, entry.Permanent + entry.Temporary})
		}
	}

	util.PrintTable(header, rows)

}

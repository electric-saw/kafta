package broker

import (
	"os"
	"sort"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	"github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

type describeBrokerOptions struct {
	config         *configuration.Configuration
	brokerIdOrAddr string
	logSize        bool
}

func NewCmdDescribeBroker(config *configuration.Configuration) *cobra.Command {
	options := &describeBrokerOptions{config: config}
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Descrive broker",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return ValidBrokers(config, len(args) > 0)
		},
		Run: func(cmd *cobra.Command, args []string) {
			options.brokerIdOrAddr = "1"
			util.CheckErr(options.complete(cmd))
			options.run()
		},
	}

	cmd.Flags().BoolVarP(&options.logSize, "log-size", "l", false, "Show log size per topic")
	return cmd

}

func (o *describeBrokerOptions) complete(cmd *cobra.Command) error {
	args := cmd.Flags().Args()
	if len(args) > 1 {
		return util.HelpErrorf(cmd, "Unexpected args: %v", args[1:])
	}
	if len(args) == 1 {
		o.brokerIdOrAddr = args[0]
	}
	return nil
}

func (d *describeBrokerOptions) run() {
	conn := kafka.MakeConnection(d.config)

	defer conn.Close()
	broker := kafka.DescribeBroker(conn, d.brokerIdOrAddr)

	sort.Strings(broker.ConsumerGroups)
	d.printConsumerGroups(broker)
	if d.logSize {
		d.printLogSize(broker)
	}

}

func (d *describeBrokerOptions) printConsumerGroups(broker *kafka.BrokerMetadata) {
	out := table.NewWriter()
	out.SetOutputMirror(os.Stdout)
	out.SetStyle(table.StyleRounded)
	out.Style().Options.SeparateRows = true
	defer out.Render()
	out.AppendHeader(table.Row{"BROKER DETAILS"})

	connected, _ := broker.Details.Broker.Connected()

	out.AppendRow(table.Row{"Id", broker.Details.Id})
	out.AppendRow(table.Row{"Host", broker.Details.Host})
	out.AppendRow(table.Row{"Address", broker.Details.Address})
	out.AppendRow(table.Row{"Connected", connected})
	out.AppendRow(table.Row{"Rack", broker.Details.Broker.Rack()})
	out.AppendRow(table.Row{"IsController", broker.Details.IsController})
}

func (d *describeBrokerOptions) printLogSize(broker *kafka.BrokerMetadata) {
	out := table.NewWriter()
	out.SetOutputMirror(os.Stdout)
	out.SetStyle(table.StyleRounded)
	out.Style().Options.SeparateRows = true

	defer out.Render()
	for _, log := range broker.Logs {
		out.AppendHeader(table.Row{"TOPIC", "PERMANENT", "TEMPORARY"})

		totalPermanent := uint64(0)
		totalTemporary := uint64(0)
		// l.AddItem(logFile.Path)
		sorted := log.SortByPermanentSize()
		for _, tLogs := range sorted {
			row := table.Row{tLogs.Topic,
				humanize.Bytes(tLogs.Permanent),
				humanize.Bytes(tLogs.Temporary),
			}

			totalPermanent += tLogs.Permanent
			totalTemporary += tLogs.Temporary

			out.AppendRow(row)
		}

		out.AppendRow(nil)

		row := table.Row{"Total",
			humanize.Bytes(totalPermanent),
			humanize.Bytes(totalTemporary),
		}

		out.AppendRow(row)
	}
}

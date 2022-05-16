package broker

import (
	"fmt"
	"os"
	"sort"

	"github.com/dustin/go-humanize"
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	"github.com/electric-saw/kafta/pkg/cmd/util"
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
	d.printDetails(broker)
	if d.logSize {
		d.printLogSize(broker)
	}

}

func (d *describeBrokerOptions) printDetails(broker *kafka.BrokerMetadata) {
	out := util.GetNewTabWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprintln(out, "BROKER DETAILS")

	connected, _ := broker.Details.Broker.Connected()

	fmt.Fprintf(out, "ID\t%d\n", broker.Details.Id)
	fmt.Fprintf(out, "Host\t%s\n", broker.Details.Host)
	fmt.Fprintf(out, "Address\t%s\n", broker.Details.Address)
	fmt.Fprintf(out, "Connected\t%v\n", connected)
	fmt.Fprintf(out, "Rack\t%s\n", broker.Details.Broker.Rack())
	fmt.Fprintf(out, "IsController\t%v\n", broker.Details.IsController)

}

func (d *describeBrokerOptions) printLogSize(broker *kafka.BrokerMetadata) {
	out := util.GetNewTabWriter(os.Stdout)
	defer out.Flush()
	fmt.Printf("sadadsa\n")
	for _, log := range broker.Logs {
		fmt.Fprintln(out, "TOPIC\tPERMANENT\tTEMPORARY")

		totalPermanent := uint64(0)
		totalTemporary := uint64(0)
		// l.AddItem(logFile.Path)
		sorted := log.SortByPermanentSize()
		for _, tLogs := range sorted {
			totalPermanent += tLogs.Permanent
			totalTemporary += tLogs.Temporary

			fmt.Fprintf(out, "%s\t%s\t%s\n", tLogs.Topic, humanize.Bytes(tLogs.Permanent), humanize.Bytes(tLogs.Temporary))
		}

		fmt.Fprintf(out, "\nTOTAL\t%s\t%s\n", humanize.Bytes(totalPermanent), humanize.Bytes(totalTemporary))

	}

}

package consumer

import (
	"os"
	"sort"

	"github.com/Shopify/sarama"
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	"github.com/jedib0t/go-pretty/table"
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
	out := table.NewWriter()
	out.SetOutputMirror(os.Stdout)
	out.SetStyle(table.StyleRounded)
	out.Style().Options.SeparateRows = true

	conn := kafka.MakeConnection(o.config)
	defer conn.Close()
	consumers := kafka.ListConsumerGroupDescriptions(conn)

	out.AppendHeader(table.Row{"NAME", "TYPE", "STATE"})

	keys := make([]string, 0, len(consumers))
	for k := range consumers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		consumerType := consumers[name]
		o.printContext(name, consumerType, out)
	}

	out.Render()

}

func (o *listConsumerOptions) printContext(name string, consumerType *sarama.GroupDescription, w table.Writer) {
	w.AppendRow(table.Row{name, consumerType.ProtocolType, consumerType.State})
}

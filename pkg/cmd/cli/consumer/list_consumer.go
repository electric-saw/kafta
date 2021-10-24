package consumer

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/Shopify/sarama"
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	"github.com/electric-saw/kafta/pkg/cmd/util"
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
	out := util.GetNewTabWriter(os.Stdout)

	conn := kafka.MakeConnection(o.config)
	defer conn.Close()
	consumers := kafka.ListConsumerGroupDescriptions(conn)

	fmt.Fprintln(out, "NAME\tTYPE\tSTATE")

	keys := make([]string, 0, len(consumers))
	for k := range consumers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		consumerType := consumers[name]
		o.printContext(name, consumerType, out)
	}

	out.Flush()

}

func (o *listConsumerOptions) printContext(name string, consumerType *sarama.GroupDescription, w io.Writer) {
	fmt.Fprintf(w, "%s\t%s\t%s\n", name, consumerType.ProtocolType, consumerType.State)
}

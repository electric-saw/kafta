package topic

import (
	"os"
	"sort"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

type listTopicOptions struct {
	config   *configuration.Configuration
	internal bool
}

func NewCmdListTopic(config *configuration.Configuration) *cobra.Command {
	options := &listTopicOptions{config: config}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List topics",
		Run: func(cmd *cobra.Command, args []string) {
			options.run()
		},
	}
	cmd.Flags().BoolVarP(&options.internal, "internal", "i", false, "Show internal topics")

	return cmd
}

func (o *listTopicOptions) run() {
	out := table.NewWriter()
	out.SetOutputMirror(os.Stdout)
	out.SetStyle(table.StyleRounded)
	out.Style().Options.SeparateRows = true

	conn := kafka.MakeConnection(o.config)
	defer conn.Close()
	topics := kafka.ListAllTopics(conn)

	out.AppendHeader(table.Row{"NAME", "PARTITIONS", "RF"})

	keys := make([]string, 0, len(topics))
	for k := range topics {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		if !o.internal && strings.HasPrefix(name, "_") {
			continue
		}

		topic := topics[name]
		o.printContext(name, topic, out)
	}

	out.Render()

}

func (o *listTopicOptions) printContext(name string, topic sarama.TopicDetail, w table.Writer) {
	w.AppendRow(table.Row{name, topic.NumPartitions, topic.ReplicationFactor})
}

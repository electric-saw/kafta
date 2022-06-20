package topic

import (
	"sort"
	"strings"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	"github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/jedib0t/go-pretty/v6/table"
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
	conn := kafka.MakeConnection(o.config)
	defer conn.Close()
	topics := kafka.ListAllTopics(conn)
	rows := []table.Row{}

	keys := make([]string, 0, len(topics))
	for k := range topics {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		if !o.internal && strings.HasPrefix(name, "_") {
			continue
		}

		rows = append(rows, table.Row{name, topics[name].NumPartitions, topics[name].ReplicationFactor})
	}

	util.PrintTable(table.Row{"name", "partitions", "replication factor"}, rows)
}

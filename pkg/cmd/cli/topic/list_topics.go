package topic

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	"github.com/electric-saw/kafta/pkg/cmd/util"
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
	out := util.GetNewTabWriter(os.Stdout)

	conn := kafka.MakeConnection(o.config)
	defer conn.Close()
	topics := kafka.ListAllTopics(conn)

	fmt.Fprint(out, "Name\tPartitions\tReplication Factor\n")

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

	out.Flush()
}

func (o *listTopicOptions) printContext(name string, topic sarama.TopicDetail, w io.Writer) {
	fmt.Fprintf(w, "%s\t%d\t%d\n", name, topic.NumPartitions, topic.ReplicationFactor)
}

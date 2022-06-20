package topic

import (
	"sort"

	"github.com/Shopify/sarama"
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type describeTopicOptions struct {
	config *configuration.Configuration
	topics []string
}

func NewCmdDescribeTopic(config *configuration.Configuration) *cobra.Command {
	options := &describeTopicOptions{config: config}
	cmd := &cobra.Command{
		Use:   "describe NAME",
		Short: "Describe a topic",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return ValidTopics(config, len(args) > 0)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd))
			options.run()
		},
	}

	return cmd

}

func (o *describeTopicOptions) complete(cmd *cobra.Command) error {
	args := cmd.Flags().Args()
	if len(args) == 0 {
		return cmdutil.HelpError(cmd, "Topic not informed")
	}
	if len(args) > 1 {
		return cmdutil.HelpError(cmd, "Only 1 Topic")
	}
	o.topics = args
	return nil
}

func (o *describeTopicOptions) run() {
	conn := kafka.MakeConnection(o.config)
	defer conn.Close()
	topic := kafka.DescribeTopics(conn, o.topics)[0]

	o.printContext(topic)

	sortedPartitions := make(map[int32]*sarama.PartitionMetadata)
	var keys SortInt32

	for _, partition := range topic.Partitions {
		sortedPartitions[partition.ID] = partition
		keys = append(keys, partition.ID)
	}

	sort.Sort(&keys)

	header := table.Row{"id", "isr", "leader", "replicas", "offline replicas"}
	rows := []table.Row{}

	for _, id := range keys {
		partition := sortedPartitions[id]
		rows = append(rows, table.Row{partition.ID, partition.Isr, partition.Leader, partition.Replicas, partition.OfflineReplicas})
	}

	cmdutil.PrintTable(header, rows)
}

func (o *describeTopicOptions) printContext(topic *sarama.TopicMetadata) {
	header := table.Row{"name", "partitions", "internal"}
	rows := []table.Row{}
	rows = append(rows, table.Row{topic.Name, len(topic.Partitions), topic.IsInternal})

	cmdutil.PrintTable(header, rows)
}

type SortInt32 []int32

func (a SortInt32) Len() int           { return len(a) }
func (a SortInt32) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortInt32) Less(i, j int) bool { return a[i] < a[j] }

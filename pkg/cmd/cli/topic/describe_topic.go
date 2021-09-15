package topic

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/printers"
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
	out := printers.GetNewTabWriter(os.Stdout)
	defer out.Flush()

	conn := kafka.MakeConnection(o.config)
	defer conn.Close()
	topic := kafka.DescribeTopics(conn, o.topics)[0]

	err := o.printContextHeaders(out)
	cmdutil.CheckErr(err)

	err = o.printContext(topic, out)
	cmdutil.CheckErr(err)

	err = o.printContextHeadersPartition(out)
	cmdutil.CheckErr(err)

	sortedPartitions := make(map[int32]*sarama.PartitionMetadata)
	var keys SortInt32

	for _, partition := range topic.Partitions {
		sortedPartitions[partition.ID] = partition
		keys = append(keys, partition.ID)
	}

	sort.Sort(&keys)

	for _, id := range keys {
		err = o.printContextPartition(sortedPartitions[id], out)
	}
	cmdutil.CheckErr(err)

}

func (o *describeTopicOptions) printContextHeaders(out io.Writer) error {
	columnNames := []string{"INTERNAL", "NAME", "PARTITIONS"}
	_, err := fmt.Fprintf(out, "%s\n", strings.Join(columnNames, "\t"))
	return err
}

func (o *describeTopicOptions) printContext(topic *sarama.TopicMetadata, w io.Writer) error {
	prefix := " "
	if topic.IsInternal {
		prefix = "*"
	}
	_, err := fmt.Fprintf(w, "%s\t%s\t%d\n", prefix, topic.Name, len(topic.Partitions))
	return err
}

func (o *describeTopicOptions) printContextHeadersPartition(out io.Writer) error {
	columnNames := []string{"ID", "ISR", "LEADER", "REPLICAS", "OFFLINE REPLICAS"}
	_, err := fmt.Fprintf(out, "\n\n%s\n", strings.Join(columnNames, "\t"))
	return err
}

func (o *describeTopicOptions) printContextPartition(partition *sarama.PartitionMetadata, w io.Writer) error {
	_, err := fmt.Fprintf(w, "%d\t%v\t%d\t%v\t%v\n", partition.ID, partition.Isr, partition.Leader, partition.Replicas, partition.OfflineReplicas)
	return err
}

type SortInt32 []int32

func (a SortInt32) Len() int           { return len(a) }
func (a SortInt32) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortInt32) Less(i, j int) bool { return a[i] < a[j] }

package consumer

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/printers"
)

type lagConsumerOptions struct {
	config  *configuration.Configuration
	groups  []string
	verbose bool
}

func NewCmdLagConsumer(config *configuration.Configuration) *cobra.Command {
	options := &lagConsumerOptions{config: config, verbose: false}
	cmd := &cobra.Command{
		Use:   "lag NAME",
		Short: "Lag a consumer",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return ValidConsumers(config, len(args) > 0)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd))
			options.run()
		},
	}

	cmd.Flags().BoolVar(&options.verbose, "verbose", false, "Show lag by partition")

	return cmd

}

func (l *lagConsumerOptions) complete(cmd *cobra.Command) error {
	args := cmd.Flags().Args()
	if len(args) == 0 {
		return cmdutil.HelpError(cmd, "Consumer not informed")
	}
	if len(args) > 1 {
		return cmdutil.HelpError(cmd, "Only 1 consumer")
	}
	l.groups = args
	return nil
}

func (l *lagConsumerOptions) run() {
	out := printers.GetNewTabWriter(os.Stdout)
	defer out.Flush()

	conn := kafka.MakeConnection(l.config)
	defer conn.Close()

	consumers := kafka.ConsumerLag(conn, l.groups)

	if !l.verbose {
		l.printTotalLag(out, consumers)
	} else {
		l.printLagByPartition(out, consumers)
	}
}

func (l *lagConsumerOptions) printTotalLag(out io.Writer, consumers map[string]*kafka.ConsumerGroupOffset) {
	fmt.Fprint(out, "CONSUMER\tTOPIC\tTOTAL LAG\n")

	for name, consumer := range consumers {
		for topicName, topic := range consumer.Topics {

			fmt.Fprintf(out, "%s\t%s\t%d\n", name, topicName, topic.GetLagTopicLag())
		}
	}
}

func (l *lagConsumerOptions) printLagByPartition(out io.Writer, consumers map[string]*kafka.ConsumerGroupOffset) {
	fmt.Fprint(out, "CONSUMER\n")
	for _, consumer := range consumers {
		fmt.Fprintf(out, "%s\n", consumer.Id)
		if len(consumer.Topics) == 0 {
			fmt.Fprintf(out, "...\tis empty")
		} else {
			for topicName, topic := range consumer.Topics {

				fmt.Fprint(out, "|   TOPIC\tTOTAL LAG\n")
				fmt.Fprintf(out, "| - %s\t%d\n", topicName, topic.GetLagTopicLag())
				fmt.Fprint(out, "|   |  PARTITION\tCONSUMER OFFSET\tPARTITION OFFSET\tLAG\n")

				keys := make(int32s, 0, len(topic.Partitions))
				for id := range topic.Partitions {
					keys = append(keys, id)
				}
				sort.Sort(keys)

				for id, partition := range topic.Partitions {
					fmt.Fprintf(out, "|   | - %d\t%d\t%d\t%d\n", id, partition.Current, partition.Max, topic.GetLagPartition(id))

				}

			}
		}
	}

}

type int32s []int32

func (a int32s) Len() int           { return len(a) }
func (a int32s) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a int32s) Less(i, j int) bool { return a[i] < a[j] }

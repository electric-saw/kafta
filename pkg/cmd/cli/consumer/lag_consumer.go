package consumer

import (
	"fmt"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
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
		Short: "list consumer lag",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return ValidConsumers(config, len(args) > 0)
		},
		Run: func(cmd *cobra.Command, args []string) {
			options.complete(cmd)
			options.run()
		},
	}

	cmd.Flags().BoolVar(&options.verbose, "verbose", false, "Show lag by partition")

	return cmd
}

func (l *lagConsumerOptions) complete(cmd *cobra.Command) {
	args := cmd.Flags().Args()
	l.groups = args
}

func (l *lagConsumerOptions) run() {
	conn := kafka.EstablishKafkaConnection(l.config)
	defer conn.Close()

	if len(l.groups) == 0 {
		for group := range kafka.ListConsumerGroups(conn) {
			l.groups = append(l.groups, group)
		}
	}

	consumers := kafka.ConsumerLag(conn, l.groups)

	if !l.verbose {
		l.printTotalLag(consumers)
	} else {
		l.printLagByPartition(consumers)
	}
}

func (l *lagConsumerOptions) printTotalLag(consumers map[string]*kafka.ConsumerGroupOffset) {
	rows := []table.Row{}

	for name, consumer := range consumers {
		for topicName, topic := range consumer.Topics {
			rows = append(rows, table.Row{name, topicName, topic.GetLagTopicLag()})
		}
	}

	cmdutil.PrintTable(table.Row{"consumer", "topic", "total lag"}, rows)
}

func (l *lagConsumerOptions) printLagByPartition(consumers map[string]*kafka.ConsumerGroupOffset) {
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	tab := table.NewWriter()
	tab.AppendHeader(
		table.Row{"consumer", "topic", "partition", "consumer offset", "partition offset", "lag"},
		rowConfigAutoMerge,
	)

	for _, consumer := range consumers {
		if len(consumer.Topics) == 0 {
			fmt.Printf("...\tis empty")
		} else {
			for topicName, topic := range consumer.Topics {
				for id, partition := range topic.Partitions {
					tab.AppendRow(table.Row{consumer.Id, topicName, id, partition.Current, partition.Max, topic.GetLagPartition(id)})
				}
			}
		}
	}

	tab.SetColumnConfigs([]table.ColumnConfig{
		{Number: 0, AutoMerge: true},
		{Number: 1, AutoMerge: true},
		{Number: 2, AutoMerge: true},
		{
			Number:      3,
			Align:       text.AlignCenter,
			AlignFooter: text.AlignCenter,
			AlignHeader: text.AlignCenter,
		},
		{
			Number:      4,
			Align:       text.AlignCenter,
			AlignFooter: text.AlignCenter,
			AlignHeader: text.AlignCenter,
		},
		{
			Number:      5,
			Align:       text.AlignCenter,
			AlignFooter: text.AlignCenter,
			AlignHeader: text.AlignCenter,
		},
		{
			Number:      6,
			Align:       text.AlignCenter,
			AlignFooter: text.AlignCenter,
			AlignHeader: text.AlignCenter,
		},
	})

	tab.SetStyle(table.StyleDefault)
	fmt.Println(tab.Render())
}

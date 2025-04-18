package consumer

import (
	"fmt"
	"sort"

	"github.com/IBM/sarama"
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type describeConsumerOptions struct {
	config *configuration.Configuration
	groups []string
}

func NewCmdDescribeConsumer(config *configuration.Configuration) *cobra.Command {
	options := &describeConsumerOptions{config: config}
	cmd := &cobra.Command{
		Use:   "describe NAME",
		Short: "Describe a consumer",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return ValidConsumers(config, len(args) > 0)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd))
			options.run()
		},
	}

	return cmd
}

func (o *describeConsumerOptions) complete(cmd *cobra.Command) error {
	args := cmd.Flags().Args()
	if len(args) == 0 {
		return cmdutil.HelpError(cmd, "Consumer not informed")
	}
	if len(args) > 1 {
		return cmdutil.HelpError(cmd, "Only 1 consumer")
	}
	o.groups = args
	return nil
}

func (o *describeConsumerOptions) run() {
	conn := kafka.EstablishKafkaConnection(o.config)
	defer conn.Close()
	consumer := kafka.DescribeConsumerGroups(conn, o.groups)[0]

	o.printBasicInfo(consumer)
	o.printMembers(consumer.Members)
}

func (o *describeConsumerOptions) printBasicInfo(group *sarama.GroupDescription) {
	header := table.Row{"id", "protocol", "protocol type", "state", "member count"}
	rows := []table.Row{}
	rows = append(
		rows,
		table.Row{
			group.GroupId,
			group.Protocol,
			group.ProtocolType,
			group.State,
			len(group.Members),
		},
	)
	cmdutil.PrintTable(header, rows)
}

func (o *describeConsumerOptions) printMembers(members map[string]*sarama.GroupMemberDescription) {
	tab := table.NewWriter()
	tab.SetStyle(table.StyleDefault)
	tab.AppendHeader(table.Row{"member id", "member host", "topic", "partitions"})

	var keys []string

	for key := range members {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, id := range keys {
		err := o.printContextPartition(members[id], tab)
		cmdutil.CheckErr(err)
	}

	tab.SetStyle(table.StyleDefault)
	fmt.Println(tab.Render())
}

func (o *describeConsumerOptions) printContextPartition(
	member *sarama.GroupMemberDescription,
	tab table.Writer,
) error {
	memberAssignment, err := member.GetMemberAssignment()
	cmdutil.CheckErr(err)
	memberMetadata, err := member.GetMemberMetadata()
	if err != nil {
		return err
	}

	if memberMetadata == nil {
		return nil
	}

	for _, topic := range memberMetadata.Topics {
		tab.AppendRow(
			table.Row{member.ClientId, member.ClientHost, topic, memberAssignment.Topics[topic]},
		)
	}
	return err
}

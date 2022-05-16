package consumer

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
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
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
	out := util.GetNewTabWriter(os.Stdout)
	defer out.Flush()

	conn := kafka.MakeConnection(o.config)
	defer conn.Close()
	consumer := kafka.DescribeConsumerGroups(conn, o.groups)[0]

	o.printBasicInfo(consumer, out)
	o.printMembers(consumer.Members, out)
}

func (o *describeConsumerOptions) printBasicInfo(group *sarama.GroupDescription, out io.Writer) {
	err := o.printTopicHeaders(out)
	cmdutil.CheckErr(err)
	err = o.printTopic(group, out)
	cmdutil.CheckErr(err)
}

func (o *describeConsumerOptions) printMembers(members map[string]*sarama.GroupMemberDescription, out io.Writer) {

	err := o.printContextHeadersPartition(out)
	cmdutil.CheckErr(err)

	var keys []string

	for key := range members {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, id := range keys {
		err = o.printContextPartition(members[id], out)
		cmdutil.CheckErr(err)
	}

}

func (o *describeConsumerOptions) printTopicHeaders(out io.Writer) error {
	columnNames := []string{"ID", "PROTOCOL", "PROTOCOL TYPE", "STATE", "MEMBER COUNT"}
	_, err := fmt.Fprintf(out, "%s\n", strings.Join(columnNames, "\t"))
	return err
}

func (o *describeConsumerOptions) printTopic(group *sarama.GroupDescription, w io.Writer) error {
	_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n", group.GroupId, group.Protocol, group.ProtocolType, group.State, len(group.Members))
	return err
}

func (o *describeConsumerOptions) printContextHeadersPartition(out io.Writer) error {
	columnNames := []string{"MEMBER ID", "MEMBER HOST", "TOPIC", "PARTITIONS"}
	_, err := fmt.Fprintf(out, "\n\n%s\n", strings.Join(columnNames, "\t"))
	return err
}

func (o *describeConsumerOptions) printContextPartition(member *sarama.GroupMemberDescription, w io.Writer) error {
	memberAssignment, err := member.GetMemberAssignment()
	cmdutil.CheckErr(err)
	memberMetadata, _ := member.GetMemberMetadata()
	for _, topic := range memberMetadata.Topics {

		_, err = fmt.Fprintf(w, "%s\t%v\t%s\t%v\n", member.ClientId, member.ClientHost, topic, memberAssignment.Topics[topic])
	}
	return err
}

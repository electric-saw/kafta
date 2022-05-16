package topic

import (
	"strings"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
)

type createTopicOptions struct {
	config       *configuration.Configuration
	name         string
	rf           int16
	partitions   int32
	topicConfigs string
}

func NewCmdCreateTopic(config *configuration.Configuration) *cobra.Command {
	options := &createTopicOptions{config: config}
	cmd := &cobra.Command{
		Use:   "create NAME [--partitions=10] [--rf=3]",
		Short: "Create topics",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd))
			cmdutil.CheckErr(options.run())
		},
	}

	cmd.Flags().Int32VarP(&options.partitions, "partitions", "p", 10, "Number of partitions")
	cmd.Flags().Int16Var(&options.rf, "rf", 3, "Number of replication on partition")
	cmd.Flags().StringVarP(&options.name, "configs", "c", "", "Configs")

	return cmd

}

func (o *createTopicOptions) complete(cmd *cobra.Command) error {
	args := cmd.Flags().Args()
	if len(args) > 1 {
		return cmdutil.HelpErrorf(cmd, "Unexpected args: %v", args)
	}
	if len(args) == 1 {
		o.name = args[0]
	}
	return nil
}

func (o *createTopicOptions) run() error {
	conn := kafka.MakeConnection(o.config)
	defer conn.Close()
	return kafka.CreateTopic(conn, o.name, o.partitions, o.rf, stringToMapPointer(o.topicConfigs))
}

// mapToMapPointer split string=string to a map[string]string
func stringToMapPointer(s string) map[string]*string {
	m := make(map[string]*string)
	for _, v := range strings.Split(s, ",") {
		kv := strings.Split(v, "=")
		if len(kv) == 2 {
			m[kv[0]] = &kv[1]
		}
	}
	return m
}

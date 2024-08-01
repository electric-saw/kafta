package topic

import (
	"fmt"
	"sort"
	"strings"

	"github.com/IBM/sarama"
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	"github.com/electric-saw/kafta/pkg/cmd/util"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type listConfigsOptions struct {
	config *configuration.Configuration
	topic  string
}

func NewCmdListConfigs(config *configuration.Configuration) *cobra.Command {
	options := &listConfigsOptions{config: config}
	cmd := &cobra.Command{
		Use:   "list-configs TOPIC",
		Short: "List all configurations for a topic",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			options.topic = args[0]
			cmdutil.CheckErr(options.run())
		},
	}

	return cmd
}

func (o *listConfigsOptions) run() error {
	conn := kafka.MakeConnection(o.config)
	defer conn.Close()

	resource := sarama.ConfigResource{
		Name: o.topic,
		Type: sarama.TopicResource,
	}
	configs, err := conn.Admin.DescribeConfig(resource)
	if err != nil {
		return fmt.Errorf("failed to describe config for topic %s: %w", o.topic, err)
	}

	rows := []table.Row{}
	for _, config := range configs {
		value := strings.TrimSpace(config.Value)
		if value == "" {
			value = "N/A"
		}

		rows = append(rows, table.Row{config.Name, value, config.Source.String()})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0].(string) < rows[j][0].(string)
	})

	util.PrintTable(
		table.Row{"Config Name", "Config Value", "Source"},
		rows,
	)

	return nil
}

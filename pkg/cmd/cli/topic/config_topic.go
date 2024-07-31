package topic

import (
	"fmt"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/olekukonko/tablewriter"
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Config Name", "Config Value", "Source"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold})
	table.SetColumnColor(tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{})

	for _, config := range configs {
		value := strings.TrimSpace(config.Value)
		if value == "" {
			value = "N/A"
		}

		table.Append([]string{config.Name, value, config.Source.String()})
	}

	table.SetAutoFormatHeaders(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(30)

	table.Render()

	return nil
}

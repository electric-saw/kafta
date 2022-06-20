package console

import (
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	consumer "github.com/electric-saw/kafta/pkg/cmd/cli/consumer"
	producer "github.com/electric-saw/kafta/pkg/cmd/cli/producer"
	"github.com/spf13/cobra"
)

func NewCmdConsole(config *configuration.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "console",
		Short: "Console management",
	}
	cmd.AddCommand(producer.NewCmdProduceMessage(config))
	cmd.AddCommand(consumer.NewCmdConsumeMessage(config))
	return cmd
}

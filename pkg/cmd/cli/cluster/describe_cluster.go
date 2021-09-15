package cluster

import (
	"os"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

type describeClusters struct {
	config *configuration.Configuration
}

func NewCmdDescribeCluster(config *configuration.Configuration) *cobra.Command {
	options := &describeClusters{config: config}
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe current-cluster",
		Run: func(cmd *cobra.Command, args []string) {
			options.run()
		},
	}
	return cmd
}

func (o *describeClusters) run() {
	out := table.NewWriter()
	out.SetOutputMirror(os.Stdout)
	out.SetStyle(table.StyleRounded)
	out.Style().Options.SeparateRows = true

	conn := kafka.MakeConnection(o.config)
	defer conn.Close()

	brokers := kafka.GetBrokers(conn)
	out.AppendHeader(table.Row{"ID", "ADDR", "CONTROLLER"})

	for _, broker := range brokers {
		o.printContext(broker.ID(), broker.Address, broker.IsController, out)
	}

	out.Render()
}

func (o *describeClusters) printContext(id int32, addr string, isController bool, w table.Writer) {
	controllerFlag := ""
	if isController {
		controllerFlag = "*"
	}
	w.AppendRow(table.Row{id, addr, controllerFlag})
}

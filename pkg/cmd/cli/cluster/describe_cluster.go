package cluster

import (
	"fmt"
	"io"
	"os"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	"github.com/electric-saw/kafta/pkg/cmd/util"
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
	out := util.GetNewTabWriter(os.Stdout)

	conn := kafka.MakeConnection(o.config)
	defer conn.Close()

	brokers := kafka.GetBrokers(conn)
	fmt.Fprintln(out, "ID\tADDR\tCONTROLLER")

	for _, broker := range brokers {
		o.printContext(broker.ID(), broker.Address, broker.IsController, out)
	}

	out.Flush()
}

func (o *describeClusters) printContext(id int32, addr string, isController bool, w io.Writer) {
	controllerFlag := ""
	if isController {
		controllerFlag = "*"
	}
	fmt.Fprintf(w, "%d\t%s\t%s\n", id, addr, controllerFlag)
}

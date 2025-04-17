package topic

import (
	"fmt"
	"strings"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
)

type resetTopicOptions struct {
	config *configuration.Configuration
	name   string
	props  []string
}

func NewCmdConfigResetTopic(config *configuration.Configuration) *cobra.Command {
	options := &resetTopicOptions{config: config}
	cmd := &cobra.Command{
		Use:   "reset NAME CONFIG [CONFIG&CONFIG2 ...]",
		Short: "reset topics",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd))
			cmdutil.CheckErr(options.run())
		},
	}

	return cmd
}

func (o *resetTopicOptions) complete(cmd *cobra.Command) error {
	args := cmd.Flags().Args()
	switch {
	case len(args) == 2:
		o.name = args[0]
		o.props = append(o.props, strings.Split(args[1], "&")...)
	default:
		return cmdutil.HelpErrorf(cmd, "Unexpected args: %v", args)
	}

	return nil
}

func (o *resetTopicOptions) run() error {
	conn := kafka.EstablishKafkaConnection(o.config)
	defer conn.Close()

	if err := kafka.ResetProp(conn, o.name, o.props); err != nil {
		return err
	} else {
		fmt.Printf("Topic %s reseted\n", o.name)
	}

	return nil
}

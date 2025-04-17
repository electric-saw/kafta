package topic

import (
	"fmt"
	"strings"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
)

type updateTopicOptions struct {
	config *configuration.Configuration
	name   string
	props  map[string]string
}

func NewCmdConfigUpdateTopic(config *configuration.Configuration) *cobra.Command {
	options := &updateTopicOptions{config: config}
	cmd := &cobra.Command{
		Use:   "update NAME CONFIG=VALUE [CONFIG=VALUE&CONFIG2=VALUE2 ...]",
		Short: "update topics",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd))
			cmdutil.CheckErr(options.run())
		},
	}

	return cmd
}

func (o *updateTopicOptions) complete(cmd *cobra.Command) error {
	args := cmd.Flags().Args()

	haveTwoArgs := len(args) == 2
	containsEqualsSign := strings.Contains(args[1], "=")

	if haveTwoArgs && containsEqualsSign {
		o.name = args[0]
		o.props = make(map[string]string)

		for prop := range strings.SplitSeq(args[1], "&") {
			keyValue := strings.Split(prop, "=")
			o.props[keyValue[0]] = keyValue[1]
		}

		return nil
	}

	return cmdutil.HelpErrorf(cmd, "Unexpected args: %v", args)
}

func (o *updateTopicOptions) run() error {
	conn := kafka.EstablishKafkaConnection(o.config)
	defer conn.Close()

	if err := kafka.SetProp(conn, o.name, o.props); err != nil {
		return err
	} else {
		fmt.Printf("Topic %s updated\n", o.name)
	}

	return nil
}

package topic

import (
	"fmt"
	"strings"

	"github.com/Songmu/prompter"
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
)

type resizeTopicOptions struct {
	config *configuration.Configuration
	name   string
	props  map[string]string
	quiet  bool
}

func NewCmdConfigResizeTopic(config *configuration.Configuration) *cobra.Command {
	options := &resizeTopicOptions{config: config}
	cmd := &cobra.Command{
		Use:   "resize NAME PARTITIONS [--quiet]",
		Short: "resize topics",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd))
			cmdutil.CheckErr(options.run())
		},
	}
	cmd.Flags().BoolVarP(&options.quiet, "quiet", "q", false, "Quiet mode")
	return cmd
}

func (o *resizeTopicOptions) complete(cmd *cobra.Command) error {
	args := cmd.Flags().Args()
	if len(args) != 2 || !strings.HasPrefix(args[1], "partitions=") {
		return cmdutil.HelpErrorf(cmd, "Unexpected args: %v", args)
	}

	o.name = args[0]
	o.props = make(map[string]string)

	keyValue := strings.SplitN(args[1], "=", 2)
	if len(keyValue) != 2 {
		return cmdutil.HelpErrorf(cmd, "Invalid partition value: %v", args[1])
	}
	o.props[keyValue[0]] = keyValue[1]

	return nil
}

func (o *resizeTopicOptions) run() error {
	if !o.quiet {
		if !prompter.YN("you really want to resize ? ", false) {
			return nil
		}
	}
	conn := kafka.MakeConnection(o.config)
	defer conn.Close()

	err := kafka.UpdatePartitions(conn, o.name, o.props)
	if err != nil {
		return err
	}

	fmt.Printf("Topic %s resized\n", o.name)

	return nil
}

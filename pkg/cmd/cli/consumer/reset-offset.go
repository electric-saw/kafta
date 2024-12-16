package consumer

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Songmu/prompter"
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
)

type resetOffsetOptions struct {
	config    *configuration.Configuration
	group     string
	topic     string
	partition int32
	offset    int64
	timestamp int64
	useOffset bool
	quiet     bool
}

func NewCmdResetOffset(config *configuration.Configuration) *cobra.Command {
	options := &resetOffsetOptions{config: config}

	cmd := &cobra.Command{
		Use:   "reset-offset GROUP TOPIC PARTITION [--offset OFFSET | --timestamp TIMESTAMP] [--quiet]",
		Short: "Reset the offset for a consumer group",
		Long: `Reset the offset for a specific group, topic, and partition.
You can reset to a specific offset or the offset corresponding to a timestamp.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd, args))
			cmdutil.CheckErr(options.run())
		},
	}

	cmd.Flags().Int64Var(&options.offset, "offset", -1, "Target offset to reset to")
	cmd.Flags().Int64Var(&options.timestamp, "timestamp", -1, "Timestamp to calculate the offset (e.g., '2024-12-01T15:04:05Z')")
	cmd.Flags().BoolVarP(&options.quiet, "quiet", "q", false, "Quiet mode")
	cmd.Flags().Int32Var(&options.partition, "partition", -1, "Partition to reset the offset for")

	return cmd
}

func (o *resetOffsetOptions) complete(cmd *cobra.Command, args []string) error {
	if len(args) != 3 {
		return cmdutil.HelpErrorf(cmd, "Invalid number of arguments. Expected GROUP, TOPIC, and PARTITION.")
	}

	o.group = args[0]
	o.topic = args[1]

	partition, err := strconv.Atoi(args[2])
	if err != nil || partition < 0 {
		return cmdutil.HelpErrorf(cmd, "Invalid partition: %s", args[2])
	}
	o.partition = int32(partition)

	switch {
	case o.offset != -1:
		o.useOffset = true
	case o.timestamp != -1:
		o.useOffset = false
		_, err := time.Parse(time.RFC3339, time.Unix(0, o.timestamp*int64(time.Millisecond)).Format(time.RFC3339))
		if err != nil {
			return cmdutil.HelpErrorf(cmd, "Invalid timestamp format: %d. Use RFC3339 format (e.g., '2024-12-01T15:04:05Z').", o.timestamp)
		}
	default:
		return cmdutil.HelpErrorf(cmd, "You must specify either --offset or --timestamp.")
	}

	return nil
}

func (o *resetOffsetOptions) run() error {
	if !o.quiet {
		message := fmt.Sprintf("Do you really want to reset the offset for group '%s', topic '%s', partition %d?", o.group, o.topic, o.partition)
		if !prompter.YN(message, false) {
			return nil
		}
	}

	conn := kafka.MakeConnection(o.config)
	defer conn.Close()

	var targetOffset int64
	var err error

	if o.useOffset {
		targetOffset = o.offset
	} else {
		timestamp := time.Unix(0, o.timestamp*int64(time.Millisecond)).Format(time.RFC3339)
		parsedTime, _ := time.Parse(time.RFC3339, timestamp)
		targetOffset, err = kafka.GetOffsetForTimestamp(conn, o.topic, o.partition, parsedTime.Unix()*int64(time.Millisecond))
		if err != nil {
			return err
		}
	}

	err = kafka.ResetConsumerGroupOffset(conn, o.group, o.topic, o.partition, targetOffset)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully reset offset for group '%s', topic '%s', partition %d to offset %d.\n",
		o.group, o.topic, o.partition, targetOffset)

	return nil
}

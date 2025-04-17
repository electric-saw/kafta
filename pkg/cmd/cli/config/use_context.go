package config

import (
	"errors"
	"fmt"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
)

const (
	useContextExample = `
		# Use the context for the kafka-dev cluster
		kafta config use-context kafka-dev`
)

type useContextOptions struct {
	kafkaconfig *configuration.KaftaConfig
	contextName string
}

func NewCmdConfigUseContext(config *configuration.Configuration) *cobra.Command {
	options := &useContextOptions{kafkaconfig: config.KaftaData}

	cmd := &cobra.Command{
		Use:                   "use-context CONTEXT_NAME",
		DisableFlagsInUseLine: true,
		Short:                 "Sets the current-context in a config file",
		Aliases:               []string{"use"},
		Long:                  `Sets the current-context in a config file`,
		Example:               useContextExample,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return ValidContexts(config, len(args) > 0)
		},

		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd))
			cmdutil.CheckErr(options.run())
			fmt.Printf("Switched to context %q.\n", options.contextName)
		},
	}

	return cmd
}

func (o *useContextOptions) run() error {
	err := o.validate()
	if err != nil {
		return err
	}
	o.kafkaconfig.CurrentContext = o.contextName

	return nil
}

func (o *useContextOptions) complete(cmd *cobra.Command) error {
	endingArgs := cmd.Flags().Args()
	if len(endingArgs) != 1 {
		return cmdutil.HelpErrorf(cmd, "Unexpected args: %v", endingArgs)
	}

	o.contextName = endingArgs[0]
	return nil
}

func (o *useContextOptions) validate() error {
	if len(o.contextName) == 0 {
		return errors.New("empty context names are not allowed")
	}

	for name := range o.kafkaconfig.Contexts {
		if name == o.contextName {
			return nil
		}
	}

	return fmt.Errorf("no context exists with the name: %q", o.contextName)
}

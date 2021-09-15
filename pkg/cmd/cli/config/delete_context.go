package config

import (
	"fmt"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
)

const (
	deleteContextExample = `
		# Delete the context for the kafka-dev cluster
		kafta config delete-context kafka-dev`
)

func NewCmdConfigDeleteContext(config *configuration.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete-context NAME",
		DisableFlagsInUseLine: true,
		Short:                 "Delete the specified context from the config",
		Example:               deleteContextExample,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return ValidContexts(config, len(args) > 0)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(runDeleteContext(config, cmd))
		},
	}

	return cmd
}

func runDeleteContext(config *configuration.Configuration, cmd *cobra.Command) error {
	args := cmd.Flags().Args()
	if len(args) != 1 {
		err := cmd.Help()
		cmdutil.CheckErr(err)
		return nil
	}

	name := args[0]
	_, ok := config.KaftaData.Contexts[name]
	if !ok {
		return fmt.Errorf("cannot delete context %s, not in %s", name, config.KaftaData.ConfigPath())
	}

	if config.KaftaData.CurrentContext == name {
		fmt.Printf("warning: this removed your active context, use \"kata config use-context\" to select a different one\n")
		config.KaftaData.CurrentContext = ""
	}

	delete(config.KaftaData.Contexts, name)

	fmt.Printf("deleted context %s from %s\n", name, config.KaftaData.ConfigPath())

	return nil
}

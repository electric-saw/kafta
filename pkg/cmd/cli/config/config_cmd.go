package config

import (
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/spf13/cobra"
)

func NewCmdConfig(config *configuration.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "config SUBCOMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Modify config files",
	}

	cmd.AddCommand(NewCmdConfigCurrentContext(config))
	cmd.AddCommand(NewCmdConfigDeleteContext(config))
	cmd.AddCommand(NewCmdConfigGetContexts(config))
	cmd.AddCommand(NewCmdConfigRenameContext(config))
	cmd.AddCommand(NewCmdConfigSetContext(config))
	cmd.AddCommand(NewCmdConfigUseContext(config))

	return cmd

}

func ValidContexts(config *configuration.Configuration, hasArgs bool) ([]string, cobra.ShellCompDirective) {
	var contexts []string
	if hasArgs {
		return contexts, cobra.ShellCompDirectiveNoFileComp
	}

	for id := range config.KaftaData.Contexts {
		contexts = append(contexts, id)
	}

	return contexts, cobra.ShellCompDirectiveNoFileComp
}

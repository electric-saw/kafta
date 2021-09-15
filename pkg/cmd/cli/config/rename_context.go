package config

import (
	"errors"
	"fmt"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
)

type RenameContextOptions struct {
	contextName string
	newName     string
}

const (
	renameContextUse = "rename-context CONTEXT_NAME NEW_NAME"

	renameContextShort = "Renames a context from the config file."

	renameContextLong = `
		Renames a context from the config file.

		CONTEXT_NAME is the context name that you wish to change.

		NEW_NAME is the new name you wish to set.

		Note: In case the context being renamed is the 'current-context', this field will also be updated.`

	renameContextExample = `
		# Rename the context 'old-name' to 'new-name' in your config file
		kafta config rename-context old-name new-name`
)

func NewCmdConfigRenameContext(config *configuration.Configuration) *cobra.Command {
	options := &RenameContextOptions{}

	cmd := &cobra.Command{
		Use:                   renameContextUse,
		DisableFlagsInUseLine: true,
		Short:                 renameContextShort,
		Long:                  renameContextLong,
		Example:               renameContextExample,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return ValidContexts(config, len(args) > 0)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.Complete(cmd, args))
			cmdutil.CheckErr(options.Validate())
			cmdutil.CheckErr(options.RunRenameContext(config))
		},
	}
	return cmd
}

func (o *RenameContextOptions) Complete(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return cmdutil.HelpErrorf(cmd, "Unexpected args: %v", args)
	}

	o.contextName = args[0]
	o.newName = args[1]
	return nil
}
func (o *RenameContextOptions) Validate() error {
	if len(o.newName) == 0 {
		return errors.New("You must specify a new non-empty context name")
	}
	return nil
}

func (o *RenameContextOptions) RunRenameContext(config *configuration.Configuration) error {
	configFile := config.KaftaData.ConfigPath()

	context, exists := config.KaftaData.Contexts[o.contextName]
	if !exists {
		return fmt.Errorf("cannot rename the context %q, it's not in %s", o.contextName, configFile)
	}

	_, newExists := config.KaftaData.Contexts[o.newName]
	if newExists {
		return fmt.Errorf("cannot rename the context %q, the context %q already exists in %s", o.contextName, o.newName, configFile)
	}

	config.KaftaData.Contexts[o.newName] = context
	delete(config.KaftaData.Contexts, o.contextName)

	if config.KaftaData.CurrentContext == o.contextName {
		config.KaftaData.CurrentContext = o.newName
	}

	fmt.Printf("Context %q renamed to %q.\n", o.contextName, o.newName)
	return nil
}

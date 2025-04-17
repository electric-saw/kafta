package config

import (
	"fmt"
	"os"
	"sort"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type GetContextsOptions struct {
	contextNames []string
	config       *configuration.Configuration
}

const (
	getContextsExample = `
		# List all the contexts in your config file
		kafta config get-contexts

		# Describe one context in your config file.
		kafta config get-contexts my-context`
)

func NewCmdConfigGetContexts(config *configuration.Configuration) *cobra.Command {
	options := &GetContextsOptions{config: config}

	cmd := &cobra.Command{
		Use:                   "get-contexts",
		DisableFlagsInUseLine: true,
		Short:                 "Describe one or many contexts",
		Long:                  "Displays one or many contexts from the config file.",
		Example:               getContextsExample,
		Run: func(cmd *cobra.Command, args []string) {
			options.Complete(cmd, args)
			options.RunGetContexts()
		},
	}

	return cmd
}

func (o *GetContextsOptions) Complete(cmd *cobra.Command, args []string) {
	o.contextNames = args
}

func (o *GetContextsOptions) RunGetContexts() {
	out := util.GetNewTabWriter(os.Stdout)
	config := o.config.KaftaData
	defer func() {
		err := out.Flush()
		if err != nil {
			fmt.Printf("Error closing controller: %v\n", err)
		}
	}()

	toPrint := []string{}
	if len(o.contextNames) == 0 {
		for name := range config.Contexts {
			toPrint = append(toPrint, name)
		}
	} else {
		for _, name := range o.contextNames {
			_, ok := config.Contexts[name]
			if ok {
				toPrint = append(toPrint, name)
			} else {
				util.CheckErr(fmt.Errorf("context %v not found", name))
			}
		}
	}

	header := table.Row{"name", "cluster", "schema registry", "ksql", "current"}
	rows := []table.Row{}

	sort.Strings(toPrint)
	for _, name := range toPrint {
		rows = append(
			rows,
			table.Row{
				name,
				config.Contexts[name].BootstrapServers[0],
				config.Contexts[name].SchemaRegistry,
				config.Contexts[name].Ksql,
				config.CurrentContext == name,
			},
		)
	}

	util.PrintTable(header, rows)
}

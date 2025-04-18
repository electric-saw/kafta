package config

import (
	"errors"
	"fmt"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
)

const (
	currentContextExample = `
		# Display the current-context
		kafta config current-context`
)

func NewCmdConfigCurrentContext(config *configuration.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "current-context",
		Short:   "Displays the current-context",
		Example: currentContextExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(RunCurrentContext(config))
		},
	}

	return cmd
}

func RunCurrentContext(config *configuration.Configuration) error {
	if config.KaftaData.CurrentContext == "" {
		err := errors.New("current-context is not set")
		return err
	}

	fmt.Printf("%s\n", config.KaftaData.CurrentContext)
	return nil
}

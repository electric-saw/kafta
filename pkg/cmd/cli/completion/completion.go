package completion

import (
	"io"
	"os"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
)

var (
	completionShells = map[string]func(out io.Writer, cmd *cobra.Command) error{
		"bash":       runCompletionBash,
		"zsh":        runCompletionZsh,
		"powershell": runCompletionPowerShell,
	}
)

func NewCmdCompletion(config *configuration.Configuration) *cobra.Command {
	shells := []string{}
	for s := range completionShells {
		shells = append(shells, s)
	}

	return &cobra.Command{
		Use:                   "completion SHELL",
		Short:                 "Output shell completion code for the specified shell (bash, zsh or powershell)",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			err := runCompletion(os.Stdout, cmd, args)
			cmdutil.CheckErr(err)
		},
		ValidArgs: shells,
	}
}

func runCompletion(out io.Writer, cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmdutil.UsageErrorf(cmd, "Shell not specified.")
	}
	if len(args) > 1 {
		return cmdutil.UsageErrorf(cmd, "Too many arguments. Expected only the shell type.")
	}
	run, found := completionShells[args[0]]
	if !found {
		return cmdutil.UsageErrorf(cmd, "Unsupported shell type %q.", args[0])
	}

	return run(out, cmd.Parent())
}

func runCompletionBash(out io.Writer, kafta *cobra.Command) error {
	return kafta.GenBashCompletion(out)
}

func runCompletionZsh(out io.Writer, kafta *cobra.Command) error {
	return kafta.GenZshCompletion(out)
}

func runCompletionPowerShell(out io.Writer, kafta *cobra.Command) error {
	return kafta.GenPowerShellCompletion(out)
}

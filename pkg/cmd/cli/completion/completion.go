package completion

import (
	"io"
	"os"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
)

func NewCmdCompletion(config *configuration.Configuration) *cobra.Command {
	return &cobra.Command{
		Use:                   "completion SHELL",
		Short:                 "Output shell completion code for the specified shell (bash, zsh or powershell)",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			err := runCompletion(os.Stdout, cmd, args)
			cmdutil.CheckErr(err)
		},
		ValidArgs: []string{
			"bash",
			"zsh",
			"powershell",
			"fish",
		},
	}
}

func runCompletion(out io.Writer, cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmdutil.UsageErrorf(cmd, "Shell not specified.")
	}
	if len(args) > 1 {
		return cmdutil.UsageErrorf(cmd, "Too many arguments. Expected only the shell type.")
	}
	fn := completionShells(args[0], cmd)
	if fn == nil {
		return cmdutil.UsageErrorf(cmd, "Unsupported shell type %q.", args[0])
	}

	return fn(out)
}

func completionShells(shell string, cmd *cobra.Command) func(out io.Writer) error {
	switch shell {
	case "bash":
		return cmd.GenBashCompletion
	case "zsh":
		return cmd.GenZshCompletion
	case "powershell":
		return cmd.GenPowerShellCompletion
	case "fish":
		return func(out io.Writer) error {
			return cmd.GenFishCompletion(out, true)
		}
	default:
		return nil
	}
}

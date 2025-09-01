package schema

import (
	"fmt"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/schema"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func NewCmdVersions(config *configuration.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "versions <subject>",
		Short: "List schema versions",
		Long:  "List all versions of a schema subject with compatibility information",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("missing required argument: subject")
			}
			if len(args) > 1 {
				return fmt.Errorf("too many arguments")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			subjectName := args[0]
			
			versions, err := schema.NewSubjectVersionsWithCompatibility(config, subjectName)
			if err != nil {
				cmdutil.CheckErr(err)
			}

			header := table.Row{"VERSION", "ID", "COMPATIBILITY"}
			rows := []table.Row{}

			for _, version := range versions {
				rows = append(rows, table.Row{
					version.Version,
					version.ID,
					version.CompatibilityLevel,
				})
			}

			cmdutil.PrintTable(header, rows)
		},
	}

	return cmd
}

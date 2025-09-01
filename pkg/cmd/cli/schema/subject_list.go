package schema

import (
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/schema"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func NewCmdSubjects(config *configuration.Configuration) *cobra.Command {
	var detail bool

	cmd := &cobra.Command{
		Use:   "subjects",
		Short: "List all subjects",
		Long:  "List all subjects from Schema Registry. Use --detail for compatibility and version information.",
		Run: func(cmd *cobra.Command, args []string) {
			if detail {
				// Show detailed information with compatibility and version count
				subjects, err := schema.NewSubjectListWithCompatibility(config)
				if err != nil {
					cmdutil.CheckErr(err)
				}

				header := table.Row{"SUBJECT", "COMPATIBILITY", "VERSIONS"}
				rows := []table.Row{}

				for _, subject := range subjects {
					rows = append(rows, table.Row{
						subject.Name,
						subject.CompatibilityLevel,
						subject.VersionCount,
					})
				}

				cmdutil.PrintTable(header, rows)
			} else {
				subjects, err := schema.NewSubjectList(config)
				if err != nil {
					cmdutil.CheckErr(err)
				}

				header := table.Row{"SUBJECT"}
				rows := []table.Row{}

				for _, subject := range subjects {
					rows = append(rows, table.Row{subject})
				}

				cmdutil.PrintTable(header, rows)
			}
		},
	}

	cmd.Flags().
		BoolVar(&detail, "detail", false, "Show detailed information including compatibility levels and version counts (slower)")

	return cmd
}

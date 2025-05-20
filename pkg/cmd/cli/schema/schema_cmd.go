package schema

import (
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/spf13/cobra"
)

func NewCmdSchema(config *configuration.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Schema Registry management",
	}

	cmd.AddCommand(NewCmdSubjectList(config))
	cmd.AddCommand(NewCmdSubjectVersion(config))
	cmd.AddCommand(NewCmdSchemaList(config))
	cmd.AddCommand(NewCmdSchemaDiff(config))

	return cmd
}

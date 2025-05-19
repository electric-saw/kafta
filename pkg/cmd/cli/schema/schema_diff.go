package schema

import (
	"log"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/schema"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
)

type schemaDiff struct {
	config  *configuration.Configuration
	subject string
	version string
	schema  string
}

func NewCmdSchemaDiff(config *configuration.Configuration) *cobra.Command {
	options := &schemaDiff{config: config}
	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Compare schemas",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd))
			options.run()
		},
	}
	cmd.Flags().String("subject", "", "The name of the subject to retrieve")
	cmd.Flags().String("version", "", "The version of the subject to retrieve")
	cmd.Flags().String("schema", "", "The schema to compare against")

	return cmd
}

func (o *schemaDiff) run() {
	jsonBytes, err := schema.NewSchemaList(o.config, o.subject, o.version)
	if err != nil {
		log.Fatal(err)
	}

	prettyJSON := cmdutil.PrettyJSON([]byte(jsonBytes))
	if prettyJSON == "" {
		log.Fatal("Failed to prettify JSON")
	}

	cmdutil.DiffJSONs(prettyJSON, o.schema)
}

func (o *schemaDiff) complete(cmd *cobra.Command) error {
	args := cmd.Flags().Args()
	if len(args) > 1 {
		return cmdutil.HelpErrorf(cmd, "Unexpected args: %v", args)
	}
	if len(args) == 1 {
		o.subject = args[0]
	}

	subject, err := cmd.Flags().GetString("subject")
	if err != nil {
		return err
	}
	o.subject = subject

	version, err := cmd.Flags().GetString("version")
	if err != nil {
		return err
	}
	o.version = version

	schema, err := cmd.Flags().GetString("schema")
	if err != nil {
		return err
	}
	o.schema = schema

	return nil
}

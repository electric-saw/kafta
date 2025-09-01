package schema

import (
	"encoding/json"
	"fmt"

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
		Use:   "diff SUBJECT [flags]",
		Short: "Compare schemas",
		Long:  "Compare a subject's schema with another schema",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmdutil.HelpErrorf(cmd, "error: Subject not informed")
			}
			if len(args) > 1 {
				return cmdutil.HelpErrorf(cmd, "error: Too many arguments")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			options.subject = args[0]
			cmdutil.CheckErr(options.complete(cmd))
			options.run()
		},
	}
	cmd.Flags().String("version", "", "The version of the subject to retrieve (default: latest)")
	cmd.Flags().String("schema", "", "The schema to compare against (required)")
	cmd.MarkFlagRequired("schema")

	return cmd
}

func (o *schemaDiff) run() {
	jsonBytes, err := schema.NewSchemaList(o.config, o.subject, o.version)
	if err != nil {
		cmdutil.CheckErr(err)
	}

	var errorResponse map[string]interface{}
	if err := json.Unmarshal([]byte(jsonBytes), &errorResponse); err == nil {
		if errorCode, exists := errorResponse["error_code"]; exists {
			if message, msgExists := errorResponse["message"]; msgExists {
				cmdutil.CheckErr(fmt.Errorf("%v", message))
			} else {
				cmdutil.CheckErr(fmt.Errorf("Schema Registry error (code: %v)", errorCode))
			}
		}
	}

	prettyJSON := cmdutil.PrettyJSON([]byte(jsonBytes))
	if prettyJSON == "" {
		cmdutil.CheckErr(fmt.Errorf("Failed to prettify JSON"))
	}

	cmdutil.DiffJSONs(prettyJSON, o.schema)
}

func (o *schemaDiff) complete(cmd *cobra.Command) error {
	version, err := cmd.Flags().GetString("version")
	if err != nil {
		return err
	}
	if version != "" {
		o.version = version
	} else {
		o.version = "latest"
	}

	schema, err := cmd.Flags().GetString("schema")
	if err != nil {
		return err
	}
	o.schema = schema

	return nil
}

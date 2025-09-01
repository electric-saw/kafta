package schema

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/schema"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
)

type schemaList struct {
	config  *configuration.Configuration
	subject string
	version string
}

func NewCmdSchemaList(config *configuration.Configuration) *cobra.Command {
	options := &schemaList{config: config}
	cmd := &cobra.Command{
		Use:   "get SUBJECT [flags]",
		Short: "Get schema by subject",
		Long:  "Get schema by subject and version from Schema Registry",
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
	return cmd
}

func (o *schemaList) run() {
	jsonBytes, err := schema.NewSchemaList(o.config, o.subject, o.version)
	if err != nil {
		log.Fatal(err)
	}

	var errorResponse map[string]interface{}
	if err := json.Unmarshal([]byte(jsonBytes), &errorResponse); err == nil {
		if errorCode, exists := errorResponse["error_code"]; exists {
			if message, msgExists := errorResponse["message"]; msgExists {
				cmdutil.CheckErr(fmt.Errorf("%v", message))
			} else {
				cmdutil.CheckErr(fmt.Errorf("schema Registry error (code: %v)", errorCode))
			}
		}
	}

	prettyJSON := cmdutil.PrettyJSON([]byte(jsonBytes))
	if prettyJSON == "" {
		log.Fatal("Failed to prettify JSON")
	}

	log.SetFlags(0)
	log.Println(prettyJSON)
}

func (o *schemaList) complete(cmd *cobra.Command) error {
	version, err := cmd.Flags().GetString("version")
	if err != nil {
		return err
	}
	if version != "" {
		o.version = version
	} else {
		versionsJSON, err := schema.NewSubjectVersion(o.config, o.subject)
		if err != nil {
			return err
		}

		var versions []int
		if err := json.Unmarshal([]byte(versionsJSON), &versions); err != nil {
			return err
		}
		if len(versions) == 0 {
			return cmdutil.HelpErrorf(cmd, "error: No versions found for subject")
		}

		o.version = strconv.Itoa(versions[len(versions)-1])
	}

	return nil
}

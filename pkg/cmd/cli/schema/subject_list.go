package schema

import (
	"log"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/schema"
	"github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type subjectList struct {
	config *configuration.Configuration
}

func NewCmdSubjectList(config *configuration.Configuration) *cobra.Command {
	options := &subjectList{config: config}
	cmd := &cobra.Command{
		Use:   "subjects-list",
		Short: "sub-list",
		Run: func(cmd *cobra.Command, args []string) {
			options.run()
		},
	}
	return cmd
}

func (o *subjectList) run() {
	subjects, err := schema.NewSubjectList(o.config)
	if err != nil {
		log.Fatal(err)
	}

	header := table.Row{"name"}
	rows := []table.Row{}

	for _, name := range subjects {
		rows = append(rows, table.Row{name})
	}

	util.PrintTable(header, rows)
}

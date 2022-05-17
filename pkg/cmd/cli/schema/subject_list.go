package schema

import (
	"log"
	"os"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/schema"
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

	t := table.NewWriter()
	t.SetStyle(table.StyleDefault)

	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name"})

	for _, name := range subjects {
		t.AppendRow(table.Row{name})
	}

	t.AppendSeparator()
	t.Render()
}

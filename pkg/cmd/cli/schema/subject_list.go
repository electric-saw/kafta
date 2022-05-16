package schema

import (
	"fmt"
	"os"
	"sort"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/schema"
	"github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
)

type subjectList struct {
	config *configuration.Configuration
}

func NewCmdSubjectList(config *configuration.Configuration) *cobra.Command {
	options := &subjectList{config: config}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list subjects",
		Run: func(cmd *cobra.Command, args []string) {
			options.run()
		},
	}
	return cmd
}

func (o *subjectList) run() {
	subjects := schema.NewSubjectList(o.config)

	out := util.GetNewTabWriter(os.Stdout)
	fmt.Fprint(out, "Name\n")

	sort.Strings(subjects)

	for _, name := range subjects {
		fmt.Fprintf(out, "%s\n", name)
	}

	out.Flush()
}

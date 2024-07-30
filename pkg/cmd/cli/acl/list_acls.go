package acl

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type listAclOptions struct {
	config *configuration.Configuration
}

func NewCmdListAcl(config *configuration.Configuration) *cobra.Command {
	options := &listAclOptions{config: config}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List acls",
		Run: func(cmd *cobra.Command, args []string) {
			options.run()
		},
	}

	return cmd
}

func (o *listAclOptions) run() {
	conn := kafka.MakeConnection(o.config)
	defer conn.Close()
	acls := kafka.ListAllAcls(conn)

	for _, acl := range acls {
		o.printResourceInfo(acl)
		o.printAcls(acl.Acls)
	}
}

func (o *listAclOptions) printResourceInfo(resourceAcl sarama.ResourceAcls) {
	header := table.Row{"resource type", "resource name", "resource pattern", "acl count"}
	rows := []table.Row{}
	rows = append(rows, table.Row{resourceAcl.Resource.ResourceType.String(), resourceAcl.Resource.ResourceName, resourceAcl.Resource.ResourcePatternType.String(), len(resourceAcl.Acls)})
	cmdutil.PrintTable(header, rows)
}

func (o *listAclOptions) printAcls(acls []*sarama.Acl) {
	tab := table.NewWriter()
	tab.SetStyle(table.StyleDefault)
	tab.AppendHeader(table.Row{"principal", "host", "operation", "permission type"})

	for _, acl := range acls {
		tab.AppendRow(table.Row{acl.Principal, acl.Host, acl.Operation.String(), acl.PermissionType.String()})
	}

	tab.SetStyle(table.StyleDefault)
	fmt.Println(tab.Render())
}

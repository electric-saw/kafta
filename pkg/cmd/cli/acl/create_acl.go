package acl

import (
	"github.com/IBM/sarama"
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"
)

type createAclOptions struct {
	config *configuration.Configuration
	AclOptions
}

func NewCmdCreateAcl(config *configuration.Configuration) *cobra.Command {
	options := &createAclOptions{config: config}
	cmd := &cobra.Command{
		Use:   "create [NAME] [--type=Topic] [--principal=User:CN=principal] [--host=*] [--operation=All] [--permission=Allow]",
		Short: "Create acl",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd))
			cmdutil.CheckErr(options.run())
		},
	}

	cmd.Flags().StringVarP(&options.resource_name, "configs", "c", "*", "Configs")

	cmd.Flags().VarP(
		enumflag.New(&options.resource_type, "type", ResourceTypeMapping, enumflag.EnumCaseInsensitive),
		"type", "t",
		"resource type can be 'Topic', 'Group', 'Cluster', 'TransactionalID', 'DelegationToken'")

	cmd.Flags().StringVarP(&options.acl_principal, "principal", "p", "", "Principal")

	cmd.Flags().StringVarP(&options.acl_host, "host", "s", "*", "Host")

	cmd.Flags().VarP(
		enumflag.New(&options.acl_operation, "operation", OperationMapping, enumflag.EnumCaseInsensitive),
		"operation", "o",
		"acl operation can be 'All', 'Read', 'Write', 'Create', 'Delete', 'Alter', 'Describe', 'ClusterAction', 'DescribeConfigs', 'AlterConfigs', 'IdempotentWrite'")

	cmd.Flags().VarP(
		enumflag.NewWithoutDefault(&options.acl_permission_type, "permission", PermissionMapping, enumflag.EnumCaseInsensitive),
		"permission", "m",
		"acl permission can be 'Deny', 'Allow'")

	return cmd
}

func (o *createAclOptions) complete(cmd *cobra.Command) error {
	args := cmd.Flags().Args()
	if len(args) > 1 {
		return cmdutil.HelpErrorf(cmd, "Unexpected args: %v", args)
	}
	if len(args) == 1 {
		o.resource_name = args[0]
	}
	if o.resource_type == sarama.AclResourceUnknown {
		o.resource_type = sarama.AclResourceTopic
	}
	if o.acl_operation == sarama.AclOperationUnknown {
		o.acl_operation = sarama.AclOperationAll
	}
	if o.acl_permission_type == sarama.AclPermissionUnknown {
		o.acl_permission_type = sarama.AclPermissionAllow
	}
	return nil
}

func (o *createAclOptions) run() error {
	conn := kafka.MakeConnection(o.config)
	defer conn.Close()
	return kafka.CreateAcl(conn, o.resource_name, o.resource_type, o.acl_principal, o.acl_host, o.acl_operation, o.acl_permission_type)
}

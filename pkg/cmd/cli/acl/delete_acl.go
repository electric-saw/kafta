package acl

import (
	"github.com/IBM/sarama"
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/internal/pkg/kafka"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"
)

type deleteAclOptions struct {
	config *configuration.Configuration
	AclOptions
}

func NewCmdDeleteAcl(config *configuration.Configuration) *cobra.Command {
	options := &deleteAclOptions{config: config}
	cmd := &cobra.Command{
		Use:   "delete [NAME] [--type=Topic] [--principal=User:CN=principal] [--host=*] [--operation=All] [--permission=Allow]",
		Short: "delete acls",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd))
			cmdutil.CheckErr(options.run())
		},
	}

	cmd.Flags().StringVarP(&options.resource_name, "configs", "c", "*", "Configs")

	cmd.Flags().VarP(
		enumflag.New(&options.resource_type, "type", ResourceTypeMapping, enumflag.EnumCaseInsensitive),
		"type", "t",
		"resource type can be 'Any', 'Topic', 'Group', 'Cluster', 'TransactionalID', 'DelegationToken'")

	cmd.Flags().StringVarP(&options.acl_principal, "principal", "p", "", "Principal")

	cmd.Flags().StringVarP(&options.acl_host, "host", "s", "*", "Host")

	cmd.Flags().VarP(
		enumflag.New(&options.acl_operation, "operation", OperationMapping, enumflag.EnumCaseInsensitive),
		"operation", "o",
		"acl operation can be 'Any', 'All', 'Read', 'Write', 'Create', 'Delete', 'Alter', 'Describe', 'ClusterAction', 'DescribeConfigs', 'AlterConfigs', 'IdempotentWrite'")

	cmd.Flags().VarP(
		enumflag.NewWithoutDefault(
			&options.acl_permission_type,
			"permission",
			PermissionMapping,
			enumflag.EnumCaseInsensitive,
		),
		"permission", "m",
		"acl permission can be 'Any', 'Deny', 'Allow'")

	return cmd
}

func (o *deleteAclOptions) complete(cmd *cobra.Command) error {
	args := cmd.Flags().Args()
	if len(args) > 1 {
		return cmdutil.HelpErrorf(cmd, "Unexpected args: %v", args)
	}
	if len(args) == 1 {
		o.resource_name = args[0]
	}
	if o.resource_type == sarama.AclResourceUnknown {
		o.resource_type = sarama.AclResourceAny
	}
	if o.acl_operation == sarama.AclOperationUnknown {
		o.acl_operation = sarama.AclOperationAny
	}
	if o.acl_permission_type == sarama.AclPermissionUnknown {
		o.acl_permission_type = sarama.AclPermissionAny
	}
	return nil
}

func (o *deleteAclOptions) run() error {
	conn := kafka.EstablishKafkaConnection(o.config)
	defer conn.Close()
	return kafka.DeleteAcl(
		conn,
		o.resource_name,
		o.resource_type,
		o.acl_principal,
		o.acl_host,
		o.acl_operation,
		o.acl_permission_type,
	)
}

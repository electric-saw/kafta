package acl

import (
	"github.com/IBM/sarama"
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/spf13/cobra"
)

func NewCmdAcl(config *configuration.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acl",
		Short: "Acls management",
	}

	cmd.AddCommand(NewCmdListAcl(config))
	cmd.AddCommand(NewCmdCreateAcl(config))
	cmd.AddCommand(NewCmdDeleteAcl(config))

	return cmd
}

type AclOptions struct {
	resource_name       string
	resource_type       sarama.AclResourceType
	acl_principal       string
	acl_host            string
	acl_operation       sarama.AclOperation
	acl_permission_type sarama.AclPermissionType
}

//nolint:exhaustive,gochecknoglobals // Ignoring unknown resource types
var ResourceTypeMapping = map[sarama.AclResourceType][]string{
	sarama.AclResourceAny:             {"Any"},
	sarama.AclResourceTopic:           {"Topic"},
	sarama.AclResourceGroup:           {"Group"},
	sarama.AclResourceCluster:         {"Cluster"},
	sarama.AclResourceTransactionalID: {"TransactionalID"},
	sarama.AclResourceDelegationToken: {"DelegationToken"},
}

//nolint:exhaustive,gochecknoglobals // Ignoring unknown resource types
var OperationMapping = map[sarama.AclOperation][]string{
	sarama.AclOperationAll:             {"All"},
	sarama.AclOperationRead:            {"Read"},
	sarama.AclOperationWrite:           {"Write"},
	sarama.AclOperationCreate:          {"Create"},
	sarama.AclOperationDelete:          {"Delete"},
	sarama.AclOperationAlter:           {"Alter"},
	sarama.AclOperationDescribe:        {"Describe"},
	sarama.AclOperationClusterAction:   {"ClusterAction"},
	sarama.AclOperationDescribeConfigs: {"DescribeConfigs"},
	sarama.AclOperationAlterConfigs:    {"AlterConfigs"},
	sarama.AclOperationIdempotentWrite: {"IdempotentWrite"},
}

//nolint:exhaustive,gochecknoglobals // Ignoring unknown resource types
var PermissionMapping = map[sarama.AclPermissionType][]string{
	sarama.AclPermissionDeny:  {"Deny"},
	sarama.AclPermissionAllow: {"Allow"},
	sarama.AclPermissionAny:   {"Any"},
}

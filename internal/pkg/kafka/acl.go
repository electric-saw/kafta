package kafka

import (
	"fmt"
	"os"

	"github.com/IBM/sarama"
	"github.com/electric-saw/kafta/pkg/cmd/util"
)

func ListAllAcls(conn *KafkaConnection) []sarama.ResourceAcls {
	filter := sarama.AclFilter{
		ResourcePatternTypeFilter: sarama.AclPatternAny,
		ResourceType:              sarama.AclResourceAny,
		PermissionType:            sarama.AclPermissionAny,
		Operation:                 sarama.AclOperationAny,
	}
	acls, err := conn.Admin.ListAcls(filter)

	util.CheckErr(err)

	return acls
}

func CreateAcl(
	conn *KafkaConnection,
	resource_name string,
	resource_type sarama.AclResourceType,
	principal string,
	host string,
	operation sarama.AclOperation,
	permission_type sarama.AclPermissionType,
) error {
	if resource_name == "" {
		fmt.Println("Resource name is required")
		os.Exit(0)
	}

	resource := sarama.Resource{
		ResourceName: resource_name,
		ResourceType: resource_type,
	}

	acl := sarama.Acl{
		Principal:      principal,
		Host:           host,
		Operation:      operation,
		PermissionType: permission_type,
	}

	if err := conn.Admin.CreateACL(resource, acl); err == nil {
		fmt.Println("Acl created")
		return nil
	} else {
		return err
	}
}

func DeleteAcl(
	conn *KafkaConnection,
	name string,
	resourceType sarama.AclResourceType,
	principal string,
	host string,
	operation sarama.AclOperation,
	permissionType sarama.AclPermissionType,
) error {
	if name == "" {
		fmt.Println("Resource name is required")
		os.Exit(0)
	}

	if _, err := conn.Admin.DeleteACL(sarama.AclFilter{
		ResourceName:              &name,
		ResourceType:              resourceType,
		Principal:                 &principal,
		Host:                      &host,
		Operation:                 operation,
		PermissionType:            permissionType,
		ResourcePatternTypeFilter: sarama.AclPatternAny,
	}, false); err == nil {
		fmt.Println("Acl deleted")
		return nil
	} else {
		return err
	}
}

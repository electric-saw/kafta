package kafka

import (
	"fmt"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/electric-saw/kafta/pkg/cmd/util"
)

func DescribeTopicConfig(conn *KafkaConnection, topic string) (configs []sarama.ConfigEntry) {
	resource := sarama.ConfigResource{
		Name: topic,
		Type: sarama.TopicResource,
	}
	configs, err := conn.Admin.DescribeConfig(resource)
	util.CheckErr(err)
	return configs
}

func GetTopicProp(conn *KafkaConnection, topic, key string) (configs []sarama.ConfigEntry) {
	resource := sarama.ConfigResource{
		Name:        topic,
		Type:        sarama.TopicResource,
		ConfigNames: []string{key},
	}
	configs, err := conn.Admin.DescribeConfig(resource)
	util.CheckErr(err)
	return configs

}

func SetProp(conn *KafkaConnection, topic string, props map[string]string) error {
	configs := DescribeTopicConfig(conn, topic)

	newConfigs := map[string]*string{}

	if numPartitions, ok := props["num.partitions"]; ok {
		partitionsRequest, err := strconv.Atoi(numPartitions)
		if err != nil {
			return err
		}

		val, err := conn.Admin.DescribeTopics([]string{topic})
		if err != nil {
			return err
		}
		currentPartitions := len(val[0].Partitions)

		if partitionsRequest > currentPartitions {
			newPartitionCount64, err := strconv.ParseInt(numPartitions, 10, 32)
			if err != nil {
				return err
			}
			newPartitionCount := int32(newPartitionCount64)

			err = increasePartitions(conn, topic, newPartitionCount)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("new partition count must be greater than current partitions, current: %d, new: %d", currentPartitions, partitionsRequest)
		}
	}

	for _, config := range configs {
		if !config.Default && config.Source != sarama.SourceStaticBroker {
			val := config.Value
			newConfigs[config.Name] = &val
		}
	}

	for key := range props {
		if key == "num.partitions" {
			continue
		}
		val := props[key]
		newConfigs[key] = &val
	}

	err := conn.Admin.AlterConfig(sarama.TopicResource, topic, newConfigs, false)

	return err
}

func ResetProp(conn *KafkaConnection, topic string, props []string) error {
	configs := DescribeTopicConfig(conn, topic)

	newConfigs := map[string]*string{}

	for _, config := range configs {
		if !config.Default && config.Source != sarama.SourceStaticBroker {
			val := config.Value
			newConfigs[config.Name] = &val
		}
	}

	for _, key := range props {
		delete(newConfigs, key)
	}

	err := conn.Admin.AlterConfig(sarama.TopicResource, topic, newConfigs, false)

	return err
}

func increasePartitions(conn *KafkaConnection, topic string, numPartitions int32) error {
	configs := DescribeTopicConfig(conn, topic)

	if len(configs) == 0 {
		return fmt.Errorf("topic %s does not exist", topic)
	}

	newConfigs := map[string]*string{}

	for _, config := range configs {
		if !config.Default && config.Source != sarama.SourceStaticBroker {
			val := config.Value
			newConfigs[config.Name] = &val
		}
	}

	err := conn.Admin.CreatePartitions(topic, numPartitions, nil, false)

	return err
}

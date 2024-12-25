package kafka

import (
	"fmt"
	"strconv"

	"github.com/IBM/sarama"
)

func UpdatePartitions(conn *KafkaConnection, topic string, props map[string]string) error {

	val, err := conn.Admin.DescribeTopics([]string{topic})
	if err != nil {
		return err
	}
	currentPartitions := len(val[0].Partitions)

	partitionsRequest, err := strconv.Atoi(props["num.partitions"])
	if err != nil {
		return err
	}

	if partitionsRequest > currentPartitions {
		newPartitionCount64, err := strconv.ParseInt(props["num.partitions"], 10, 32)
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

	return nil
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

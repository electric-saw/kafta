package kafka

import (
	"fmt"
	"os"

	"github.com/electric-saw/kafta/pkg/cmd/util"

	"github.com/IBM/sarama"
)

func ListAllTopics(conn *KafkaConnection) map[string]sarama.TopicDetail {
	topics, err := conn.Admin.ListTopics()
	util.CheckErr(err)
	return topics
}

func DescribeTopics(conn *KafkaConnection, topics []string) []*sarama.TopicMetadata {
	response, err := conn.Admin.DescribeTopics(topics)
	util.CheckErr(err)

	return response
}

func CreateTopic(conn *KafkaConnection, topic string, numPartitions int32, replicationFactor int16, configs map[string]*string) error {
	if topic == "" {
		fmt.Println("Topic name is required")
		os.Exit(0)
	}

	if err := conn.Admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
		ConfigEntries:     configs,
	}, false); err == nil {
		fmt.Println("Topic created")
		return err
	} else {
		return err
	}
}

func DeleteTopic(conn *KafkaConnection, topic string) error {
	if err := conn.Admin.DeleteTopic(topic); err == nil {
		fmt.Println("Topic deleted")
		return err
	} else {
		return err
	}
}

func GetTopicOffsets(conn *KafkaConnection, topic string) map[int32]int64 {
	result := make(map[int32]int64)

	partitions, err := conn.Client.Partitions(topic)
	util.CheckErr(err)

	for _, partition := range partitions {
		offset, err := conn.Client.GetOffset(topic, partition, sarama.OffsetNewest)
		util.CheckErr(err)

		result[partition] = offset
	}
	return result
}

package kafka

import (
	"github.com/IBM/sarama"
	"github.com/electric-saw/kafta/pkg/cmd/util"
)

func ListConsumerGroups(conn *KafkaConnection) map[string]string {
	consumers, err := conn.Admin.ListConsumerGroups()
	util.CheckErr(err)
	return consumers
}

func ListConsumerGroupDescriptions(conn *KafkaConnection) map[string]*sarama.GroupDescription {
	consumers, err := conn.Admin.ListConsumerGroups()
	util.CheckErr(err)

	var groups []string

	for consumer := range consumers {
		groups = append(groups, consumer)
	}

	consumersDescriptions := DescribeConsumerGroups(conn, groups)
	result := make(map[string]*sarama.GroupDescription)

	for _, consumerDescription := range consumersDescriptions {
		result[consumerDescription.GroupId] = consumerDescription
	}

	return result
}

func DescribeConsumerGroups(conn *KafkaConnection, groups []string) []*sarama.GroupDescription {
	consumers, err := conn.Admin.DescribeConsumerGroups(groups)
	util.CheckErr(err)
	return consumers
}

// ListConsumerGroupOffsets(group string, topicPartitions map[string][]int32)

func DeleteConsumerGroup(conn *KafkaConnection, group string) error {
	return conn.Admin.DeleteConsumerGroup(group)
}

type topicOffset map[int32]int64

func ConsumerLag(conn *KafkaConnection, groups []string) map[string]*ConsumerGroupOffset {

	result := make(map[string]*ConsumerGroupOffset)
	cacheTopic := make(map[string]topicOffset)

	for _, group := range DescribeConsumerGroups(conn, groups) {
		manager, err := sarama.NewOffsetManagerFromClient(group.GroupId, conn.Client)
		util.CheckErr(err)

		groupOffset := MakeConsumerGroupOffset(group.GroupId)
		result[group.GroupId] = groupOffset

		for _, member := range group.Members {
			memberAssignment, err := member.GetMemberAssignment()
			util.CheckErr(err)

			for topic, partitions := range memberAssignment.Topics {
				assigment := groupOffset.AddTopic(topic)
				var offsetsTopic topicOffset
				if data, ok := cacheTopic[topic]; !ok {
					offsetsTopic = GetTopicOffsets(conn, topic)
					cacheTopic[topic] = offsetsTopic
				} else {
					offsetsTopic = data
				}

				for _, partition := range partitions {
					patitionManager, err := manager.ManagePartition(topic, partition)
					util.CheckErr(err)

					offset, _ := patitionManager.NextOffset()

					assigment.AddOffset(partition, offset, offsetsTopic[partition])
				}
			}
		}
	}

	return result
}

func ResetConsumerGroupOffset(conn *KafkaConnection, group string, topic string, partition int32, offset int64) error {
	manager, err := sarama.NewOffsetManagerFromClient(group, conn.Client)
	util.CheckErr(err)

	partitionManager, err := manager.ManagePartition(topic, partition)
	util.CheckErr(err)

	partitionManager.ResetOffset(offset, "")
	util.CheckErr(err)
	return nil
}

func GetOffsetForTimestamp(conn *KafkaConnection, topic string, partition int32, timestamp int64) (int64, error) {
	return conn.Client.GetOffset(topic, partition, timestamp)
}

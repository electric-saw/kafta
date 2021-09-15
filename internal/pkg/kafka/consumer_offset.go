package kafka

// Consumer
//  | - topic
//    | - partition
//      | - offset

type ConsumerGroupOffset struct {
	Id     string
	Topics map[string]*TopicAssignment
}

type TopicAssignment struct {
	Partitions map[int32]PartitionOffsets
}

type PartitionOffsets struct {
	Current int64
	Max     int64
}

func MakeConsumerGroupOffset(groupId string) *ConsumerGroupOffset {
	return &ConsumerGroupOffset{
		Id:     groupId,
		Topics: make(map[string]*TopicAssignment),
	}
}

func (c *ConsumerGroupOffset) AddTopic(name string) *TopicAssignment {
	if member, ok := c.Topics[name]; ok {
		return member
	} else {
		member := &TopicAssignment{
			Partitions: make(map[int32]PartitionOffsets),
		}

		c.Topics[name] = member
		return member
	}
}

func (m *TopicAssignment) AddOffset(partition int32, consumerOffset int64, partitionOffset int64) *PartitionOffsets {
	result := PartitionOffsets{
		Current: consumerOffset,
		Max:     partitionOffset,
	}

	m.Partitions[partition] = result
	return &result
}

func (m *TopicAssignment) GetLagPartition(partitionNo int32) int64 {
	partition := m.Partitions[partitionNo]
	if result := (partition.Max - partition.Current); result < 0 {
		return 0
	} else {
		return result
	}
}

func (m *TopicAssignment) GetLagTopicLag() int64 {
	totalLag := int64(0)
	for partition := range m.Partitions {
		totalLag += m.GetLagPartition(partition)
	}
	return totalLag
}

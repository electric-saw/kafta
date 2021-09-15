package kafka

import (
	"github.com/Shopify/sarama"
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

func SetProp(conn *KafkaConnection, topic, key, value string) {
	err := conn.Admin.AlterConfig(sarama.BrokerResource, topic, map[string]*string{
		key: &value,
	}, true)

	util.CheckErr(err)
}

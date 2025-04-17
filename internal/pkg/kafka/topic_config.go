package kafka

import (
	"github.com/IBM/sarama"
	"github.com/electric-saw/kafta/pkg/cmd/util"
)

func DescribeTopicConfig(conn *KafkaConnection, topic string) []sarama.ConfigEntry {
	resource := sarama.ConfigResource{
		Name: topic,
		Type: sarama.TopicResource,
	}
	configs, err := conn.Admin.DescribeConfig(resource)
	util.CheckErr(err)
	return configs
}

func GetTopicProp(conn *KafkaConnection, topic, key string) []sarama.ConfigEntry {
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

	for _, config := range configs {
		if !config.Default && config.Source != sarama.SourceStaticBroker {
			val := config.Value
			newConfigs[config.Name] = &val
		}
	}

	for key := range props {
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

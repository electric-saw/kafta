package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/electric-saw/kafta/pkg/cmd/util"
)

func DescribeBrokerConfig(conn *KafkaConnection, brokerId string) (configs []sarama.ConfigEntry) {
	resource := sarama.ConfigResource{
		Name: brokerId,
		Type: sarama.BrokerResource,
	}
	configs, err := conn.Admin.DescribeConfig(resource)
	util.CheckErr(err)
	return configs
}

func GetBrokerProp(conn *KafkaConnection, brokerId, key string) (configs []sarama.ConfigEntry) {
	resource := sarama.ConfigResource{
		Name:        brokerId,
		Type:        sarama.BrokerResource,
		ConfigNames: []string{key},
	}
	configs, err := conn.Admin.DescribeConfig(resource)
	util.CheckErr(err)
	return configs

}

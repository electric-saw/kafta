package kafka

import (
	"fmt"
	"sort"

	"github.com/electric-saw/kafta/pkg/cmd/util"
)

func getBrokerByIdOrAddr(conn *KafkaConnection, idOrAddr any) (*Broker, error) {
	brokers := GetBrokers(conn)

	switch idOrAddr.(type) {
	case int:
		for _, broker := range brokers {
			if broker.Id == idOrAddr {
				return broker, nil
			}
		}
	case string:
		for _, broker := range brokers {
			if broker.Host == idOrAddr {
				return broker, nil
			}
		}
	}
	return nil, fmt.Errorf("Broker %s not found", idOrAddr)
}

func GetBrokers(conn *KafkaConnection) BrokersById {
	rawBrokers := conn.Client.Brokers()
	var result BrokersById

	controller, err := conn.Client.Controller()
	util.CheckErr(err)
	defer controller.Close()

	for _, rawBroker := range rawBrokers {
		result = append(result, NewBroker(rawBroker, controller.ID()))
	}

	sort.Sort(result)

	return result
}

func DescribeBroker(conn *KafkaConnection, idOrAddr any) *BrokerMetadata {
	metadata := &BrokerMetadata{
		Logs:    make([]*LogFile, 0),
		Details: &Broker{},
	}
	broker, err := getBrokerByIdOrAddr(conn, idOrAddr)
	util.CheckErr(err)

	metadata.Details = broker

	logDirs, err := conn.Admin.DescribeLogDirs([]int32{int32(broker.Id)})
	util.CheckErr(err)

	for _, logDir := range logDirs {

		for _, dir := range logDir {
			logFile := newLogFile(dir.Path)
			for _, topic := range dir.Topics {
				for _, partition := range topic.Partitions {
					logFile.set(topic.Topic, partition.Size, partition.IsTemporary)
				}
			}
			metadata.Logs = append(metadata.Logs, logFile)
		}
	}

	err = broker.Open(conn.Client.Config())
	util.CheckErr(err)
	defer broker.Close()

	return metadata
}

func GetConfigBroker(conn *KafkaConnection, broker string) {

}

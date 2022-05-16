package kafka

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/electric-saw/kafta/pkg/cmd/util"
)

func getBrokerByIdOrAddr(conn *KafkaConnection, idOrAddr string) (*Broker, error) {
	brokers := GetBrokers(conn)
	if id, err := strconv.ParseInt(idOrAddr, 10, 64); err == nil {
		for _, broker := range brokers {
			if broker.Id == int32(id) {
				return broker, nil
			}
		}
	} else {
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

func DescribeBroker(conn *KafkaConnection, idOrAddr string) *BrokerMetadata {
	metadata := &BrokerMetadata{
		Logs:    make([]*LogFile, 0),
		Details: &Broker{},
	}
	broker, err := getBrokerByIdOrAddr(conn, idOrAddr)
	util.CheckErr(err)

	metadata.Details = broker

	err = broker.Open(conn.Client.Config())
	util.CheckErr(err)
	defer broker.Close()

	return metadata
}

func GetConfigBroker(conn *KafkaConnection, broker string) {

}

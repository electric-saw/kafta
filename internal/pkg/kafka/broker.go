package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

const ControllerBrokerLabel = "*"

type Broker struct {
	Address      string
	Id           int32
	Host         string
	IsController bool
	*sarama.Broker
}

func NewBroker(broker *sarama.Broker, controllerId int32) *Broker {
	address := broker.Addr()
	id := broker.ID()
	return &Broker{
		Address:      address,
		Host:         removePort(address),
		Id:           id,
		IsController: controllerId == id,
		Broker:       broker,
	}
}

func (b *Broker) MarkedHostName() string {
	if b.IsController {
		return b.Host + ControllerBrokerLabel
	}
	return b.Host
}

func (b *Broker) String() string {
	if b == nil {
		return ""
	}
	return fmt.Sprintf("%d/%s", b.Id, b.Host)
}

type BrokersById []*Broker

func (b BrokersById) Len() int {
	return len(b)
}

func (b BrokersById) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b BrokersById) Less(i, j int) bool {
	return b[i].Id < b[j].Id
}

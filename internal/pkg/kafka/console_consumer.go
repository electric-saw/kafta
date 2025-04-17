package kafka

import (
	"context"
	"encoding/binary"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/electric-saw/kafta/pkg/cmd/util"
)

type Consumer struct {
	ready chan bool
	conn  *KafkaConnection
}

func ConsumeMessage(conn *KafkaConnection, topic string, group string, verbose bool) error {
	conn.Client.Config().ClientID = group

	keepRunning := true
	log.Printf("Initializing Consumer with group [%s]...", group)
	if verbose {
		sarama.Logger = log.New(os.Stdout, "[kafta] ", log.LstdFlags)
	}
	consumer := Consumer{
		ready: make(chan bool),
		conn:  conn,
	}

	ctx, cancel := context.WithCancel(context.Background())

	cgConfig := conn.Client.Config()

	cgConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(conn.Context.BootstrapServers, group, cgConfig)

	util.CheckErr(err)
	return err

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(consumer.ready)
		defer cancel()
		for {
			err := client.Consume(ctx, []string{topic}, &consumer)
			util.CheckErr(err)
			if ctx.Err() != nil || err != nil {
				break
			}
			
		}
	}()

	select {
	case <-consumer.ready:
	case <-ctx.Done():
		return nil
	}
	log.Println("Consumer running, waiting for events...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(sigterm)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		}

		time.Sleep(100 * time.Millisecond)
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
	return nil
}

func (consumer *Consumer) Setup(sess sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		consumer.printMessage(message)
		session.MarkMessage(message, "")

	}
	return nil
}

func (consumer *Consumer) printMessage(message *sarama.ConsumerMessage) {
	if consumer.conn.SchemaRegistryClient != nil {
		if len(message.Value) > 0 && message.Value[0] == 0 {
			schemaID := binary.BigEndian.Uint32(message.Value[1:5])
			schema, err := consumer.conn.SchemaRegistryClient.GetSchema(int(schemaID))
			if err != nil {
				log.Printf("Error getting schema for topic [%s]: %v", message.Topic, err)
			} else {
				native, _, err := schema.Codec().NativeFromBinary(message.Value[5:])
				value, _ := schema.Codec().TextualFromNative(nil, native)
				if err != nil {
					log.Printf("Error decoding message for topic [%s]: %v", message.Topic, err)
				} else {
					log.Printf("Partition: %d Key: %s Message: %s", message.Partition, string(message.Key), string(value))
				}
			}
		}
	} else {
		log.Printf("Partition: %v Key: %s Message: %s", message.Partition, string(message.Key), string(message.Value))
	}
}

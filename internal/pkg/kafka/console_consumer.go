package kafka

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/electric-saw/kafta/pkg/cmd/util"
)

type Consumer struct {
	ready chan bool
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
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(conn.Context.BootstrapServers, group, conn.Client.Config())

	util.CheckErr(err)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			err := client.Consume(ctx, []string{topic}, &consumer)
			util.CheckErr(err)
			if ctx.Err() != nil {
				break
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready
	log.Println("Consumer running, waiting for events...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(sigterm)
	defer signal.Stop(sigusr1)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
			break
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
	return nil
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if len(message.Value) == 0 {
			message.Value = message.Key
			message.Key = nil
		}
		log.Printf("Partition: %v Key: %s Message: %s", message.Partition, string(message.Key), string(message.Value))
		session.MarkMessage(message, "")
	}
	return nil
}

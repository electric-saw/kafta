package kafka

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/electric-saw/kafta/pkg/cmd/util"
)

func ProduceMessage(conn *KafkaConnection, topic string) error {
	var messageInput string
	var key *string
	producer, err := sarama.NewAsyncProducer(conn.Config.GetContext().BootstrapServers, conn.Client.Config())
	util.CheckErr(err)
	defer producer.AsyncClose()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigterm)

producerLoop:
	for {

		consoleReader := bufio.NewReader(os.Stdin)
		fmt.Print(">")
		input, err := consoleReader.ReadString('\n')
		util.CheckErr(err)

		if input == "\\quit\n" {
			producer.AsyncClose()
			break producerLoop
		}

		inputList := strings.Split(input, ":")

		if len(inputList) > 0 {
			key = &inputList[0]
			messageInput = fmt.Sprint(strings.Join(inputList[1:], ":"))
		} else {
			key = nil
			messageInput = input
		}

		message := &sarama.ProducerMessage{Topic: topic, Key: sarama.StringEncoder(*key), Value: sarama.StringEncoder(messageInput)}
		select {
		case producer.Input() <- message:
		case <-sigterm:
			log.Println("terminating: via signal")
			producer.AsyncClose()
			break producerLoop
		}
	}
	return nil
}

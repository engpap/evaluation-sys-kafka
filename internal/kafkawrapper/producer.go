package kafkawrapper

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func CreateProducer() (*kafka.Producer, error) {
	if len(os.Args) != 2 {
		fmt.Fprint(os.Stderr, "Usage: %s <config-file-path>\n", os.Args[0])
		os.Exit(1)
	}
	configFile := os.Args[1] // this can be either getting-started.properties or client.properties
	conf := ReadConfig(configFile)

	p, err := kafka.NewProducer(&conf)
	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	go producerMessageHandler(p)

	return p, nil
}

func producerMessageHandler(p *kafka.Producer) {
	for e := range p.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
			} else {
				fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
					*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
			}
		}
	}
}

// Setup clean shutdown on Ctrl+C (SIGINT) or SIGTERM
func SetupCloseProducerHandler(producer *kafka.Producer) {
	fmt.Println("Press Ctrl+C to exit.")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigchan
		fmt.Printf("Caught signal %v: terminating\n", sig)
		producer.Close()
		fmt.Println("Kafka producer closed.")
		os.Exit(0)
	}()
}

func ProduceMessage(producer *kafka.Producer, actionType string, topic string, data any) error {
	// convert course into a byte array
	byteData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// send the message to the Kafka topic
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte("data"),
		Value:          byteData,
		Headers:        []kafka.Header{{Key: "action_type", Value: []byte(actionType)}},
	}, nil)
	return err
}

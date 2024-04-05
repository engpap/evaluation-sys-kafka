package kafka

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// To handle the message consumed from Kafka, we define a callback function
type MessageHandler func(message interface{})

func CreateConsumer(topic string, handler MessageHandler) {
	if len(os.Args) != 2 {
		fmt.Fprint(os.Stderr, "Usage: %s <config-file-path>\n", os.Args[0])
		os.Exit(1)
	}
	configFile := os.Args[1]

	conf := ReadConfig(configFile)
	conf["group.id"] = "evaluation-sys"
	conf["auto.offset.reset"] = "earliest"

	c, err := kafka.NewConsumer(&conf)
	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	c.SubscribeTopics([]string{topic}, nil)

	SetupCloseConsumerHandler(c)

	go consumerMessageHandler(c, handler)
}

func consumerMessageHandler(c *kafka.Consumer, handler MessageHandler) {
	run := true
	for run {
		ev, err := c.ReadMessage(100 * time.Millisecond)
		if err != nil {
			// Errors are informational and automatically handled by the consumer
			continue
		}
		fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
			*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
		// Convert JSON to map
		var data interface{}
		err = json.Unmarshal(ev.Value, &data)
		if err != nil {
			fmt.Printf("Failed to unmarshal data: %s", err)
		}
		// execute the callback function
		handler(data)
	}
}

// Setup clean shutdown on Ctrl+C (SIGINT) or SIGTERM
func SetupCloseConsumerHandler(consumer *kafka.Consumer) {
	fmt.Println("Press Ctrl+C to exit.")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigchan
		fmt.Printf("Caught signal %v: terminating\n", sig)
		consumer.Close()
		fmt.Println("Kafka consumer closed.")
		os.Exit(0)
	}()
}

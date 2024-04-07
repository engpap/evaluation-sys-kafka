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
type MessageHandler func(action_type string, message interface{})

// CreateConsumer creates a Kafka consumer and subscribes to a topic.
//
// CONFIGS
// The consumer group ID is set to "evaluation-sys".
// The auto.offset.reset is set to "earliest" to consume all messages from the beginning of the history.
// auto.offset.reset can be either "earliest" or "latest". However, if the group.id is the same, the two options are equivalent.
// Because the consumer group has committed offsets before shutting down, upon restarting, it resumes from those offsets.
// The setting would only apply if you were to join a new consumer group to the topic or if the committed offsets were invalid.
//
// FAULT RECOVERY
// The function resetPartitionsOffset is called to seek to the beginning of the partition.
//
// CLEAN SHUTDOWN
// The function SetupCloseConsumerHandler is called to handle clean shutdown on Ctrl+C (SIGINT) or SIGTERM.
//
// HANDLING MESSAGES
// The function consumerMessageHandler is called to consume messages from the topic.
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

	printMetadata(c) // Logs on the console the metadata of the consumer

	resetPartitionsOffset(c)

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
		// Commit
		c.CommitMessage(ev)
		fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
			*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
		// Convert JSON to map
		var data interface{}
		err = json.Unmarshal(ev.Value, &data)
		if err != nil {
			fmt.Printf("Failed to unmarshal data: %s", err)
		}
		headers := ev.Headers
		actionType := ""
		for _, header := range headers {
			if string(header.Key) == "action_type" {
				actionType = string(header.Value)
				break
			}
		}
		// check that action_type is either "add" or "delete"
		if actionType != "add" && actionType != "delete" {
			fmt.Printf("Invalid action type: %s\n", actionType)
			continue
		}
		// execute the callback function
		handler(actionType, data)
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

// This function seeks to the beginning of the partition so that we can consume all messages
// from the beginning of the history.
// This is useful for fault recovery.
// When a service crashes and it is restarted, it can consume all messages from the beginning of history.
//
// Polling is necessary to avoid having 0 assigned partitions.
// We  have to wait until the consumer has completed its group join and partition assignment process.
//
// OBSERVATION: Message recovery can also be done by setting `conf["group.id"] = fmt.Sprintf("evaluation-sys-%d", time.Now().Unix())`
// However, this approach will create a new consumer group every time the service is restarted.
func resetPartitionsOffset(c *kafka.Consumer) {
	var assignedPartitions []kafka.TopicPartition
	for {
		assignedPartitions, _ = c.Assignment()
		if len(assignedPartitions) > 0 {
			break
		}
		c.Poll(100)
	}
	// for each assigned partition, seek to the beginning of the partition
	for _, partition := range assignedPartitions {
		fmt.Printf("Topic: %s, Partition: %d\n", *partition.Topic, partition.Partition)
		c.Assign([]kafka.TopicPartition{{Topic: partition.Topic, Partition: partition.Partition, Offset: kafka.OffsetBeginning}})
	}
}

// printMetadata prints the topics that the consumer is subscribed to and the metadata of each topic
func printMetadata(c *kafka.Consumer) {
	subscribedTopics, err := c.Subscription()
	if err != nil {
		fmt.Printf("Failed to get subscription: %s", err)
	}
	for _, topic := range subscribedTopics {
		fmt.Printf("Subscribed to topic: %s\n", topic)
		metadata, err := c.GetMetadata(&topic, false, 5000)
		if err != nil {
			fmt.Printf("Failed to get metadata: %s", err)
		}
		fmt.Printf("Metadata for topic %s: %v\n", topic, metadata)
	}
}

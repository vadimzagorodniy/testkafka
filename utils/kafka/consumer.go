package fioKafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"os"
	"time"
)

func Consume(FIOjobCh chan string) {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{os.Getenv("FIO_TOPIC_NAME"), os.Getenv("FIO_FAILED_TOPIC_NAME")}, nil)

	// A signal handler or similar could be used to set this to false to break the loop.
	run := true

	for run {

		msg, err := c.ReadMessage(time.Second)

		if err == nil {

			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

			switch *msg.TopicPartition.Topic {

			case os.Getenv("FIO_TOPIC_NAME"):
				FIOjobCh <- string(msg.Value)
			case os.Getenv("FIO_FAILED_TOPIC_NAME"):
				fmt.Println("Handling failed FIO")
			}

		} else if !err.(kafka.Error).IsTimeout() {
			// The client will automatically try to recover from all errors.
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}

package fioKafka

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"testjun/controllers"
)

func Produce(KafkaSendableChannel chan controllers.KafkaSendable) {

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)

	for job := range KafkaSendableChannel {
		topic := job.GetKafkaCredentials()
		fmt.Println(topic)
		v, err := json.Marshal(job)
		stringJson := string(v)
		fmt.Println(stringJson)
		if err != nil {
			fmt.Println("Invalid JSON")
		}
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(stringJson),
		}, nil)
	}

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}

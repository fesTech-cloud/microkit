package main

import (
	"context"
	"fmt"
	"log"

	"github.com/festus/microkit/adapters/rabbitmq"
	"github.com/festus/microkit/messaging"
)

func main() {
	ctx := context.Background()
	topic := "example.topic"

	// 1. Connect to RabbitMQ
	conn, err := rabbitmq.NewConnection("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 2. Create a producer
	prodCfg := messaging.Config{}
	producer, err := rabbitmq.NewProducer(conn, prodCfg)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// 3. Publish some messages
	for i := 1; i <= 3; i++ {
		payload := fmt.Sprintf("Hello World %d", i)
		msg := messaging.Message{Payload: []byte(payload)}
		if err := producer.Publish(ctx, topic, msg); err != nil {
			log.Printf("Failed to publish message: %v", err)
		} else {
			fmt.Println("Published message:", payload)
		}
	}

	fmt.Println("Producer finished publishing messages")
}

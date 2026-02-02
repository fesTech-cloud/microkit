package main

import (
	"context"
	"fmt"
	"log"

	"github.com/festech-cloud/microkit/adapters/rabbitmq"
	"github.com/festech-cloud/microkit/messaging"
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

	// 2. Create a consumer
	consCfg := messaging.Config{}
	consumer, err := rabbitmq.NewConsumer(conn, consCfg)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// 3. Subscribe and handle messages
	err = consumer.Subscribe(ctx, topic, func(ctx context.Context, msg messaging.Message) error {
		fmt.Println("Received message:", string(msg.Payload))
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	// Keep consumer running
	fmt.Println("Consumer is listening. Press Ctrl+C to exit...")
	select {}
}

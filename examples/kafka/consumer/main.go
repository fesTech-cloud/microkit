package main

import (
	"context"
	"fmt"

	"github.com/festech-cloud/microkit/adapters/kafka"
)

func main() {
	ctx := context.Background()
	topic := "example-topic"
	groupID := "example-group"

	// 1. Create Kafka connection
	conn := kafka.NewConnection([]string{"localhost:9092"})

	// 2. Create a consumer
	consumer := kafka.NewConsumer(conn, topic, groupID)
	defer consumer.Close()

	// 3. Subscribe and handle messages
	consumer.Subscribe(ctx, func(msg []byte) error {
		fmt.Printf("Received message: %s\n", string(msg))
		return nil
	})

	// Keep consumer running
	fmt.Println("Consumer is listening. Press Ctrl+C to exit...")
	select {}
}

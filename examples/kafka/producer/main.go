package main

import (
	"context"
	"fmt"
	"log"

	"github.com/festus/microkit/adapters/kafka"
)

func main() {
	ctx := context.Background()
	topic := "example-topic"

	// 1. Create Kafka connection
	conn := kafka.NewConnection([]string{"localhost:9092"})

	// 2. Create a producer
	producer := kafka.NewProducer(conn, topic)
	defer producer.Close()

	// 3. Publish some messages
	for i := 1; i <= 3; i++ {
		key := fmt.Sprintf("key-%d", i)
		payload := fmt.Sprintf("Hello Kafka %d", i)

		if err := producer.Publish(ctx, []byte(key), []byte(payload)); err != nil {
			log.Printf("Failed to publish message: %v", err)
		} else {
			fmt.Printf("Published message: key=%s, payload=%s\n", key, payload)
		}
	}

	fmt.Println("Producer finished publishing messages")
}

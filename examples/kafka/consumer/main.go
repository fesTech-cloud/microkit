package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/festus/microkit/adapters/kafka"
	"github.com/festus/microkit/internal/retry"
)

func main() {
	ctx := context.Background()
	topic := "example-topic"
	groupID := "example-group"

	// 1. Create Kafka connection
	conn := kafka.NewConnection([]string{"localhost:9092"})

	// 2. Create consumer with retry and DLQ configuration
	config := kafka.ConsumerConfig{
		RetryConfig: retry.Config{
			MaxAttempts:  3,
			InitialDelay: 100 * time.Millisecond,
			MaxDelay:     2 * time.Second,
			Multiplier:   2.0,
		},
		EnableDLQ: true,
		DLQTopic:  "example-topic-dlq",
	}

	consumer := kafka.NewConsumerWithConfig(conn, topic, groupID, config)
	defer consumer.Close()

	// 3. Subscribe with handler that sometimes fails
	consumer.Subscribe(ctx, func(msg []byte) error {
		fmt.Printf("Processing message: %s\n", string(msg))

		// Simulate random failures for demo
		if rand.Float32() < 0.3 {
			return errors.New("simulated processing error")
		}

		fmt.Printf("Successfully processed: %s\n", string(msg))
		return nil
	})

	fmt.Println("Consumer with retry/DLQ is listening. Press Ctrl+C to exit...")
	select {}
}

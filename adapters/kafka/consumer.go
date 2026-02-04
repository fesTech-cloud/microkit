// adapters/kafka/consumer.go
package kafka

import (
	"context"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

type Consumer struct {
	conn    *Connection
	topic   string
	groupID string
	r       *kafka.Reader
}

func NewConsumer(conn *Connection, topic, groupID string) *Consumer {
	return &Consumer{
		conn:    conn,
		topic:   topic,
		groupID: groupID,
		r:       conn.Reader(topic, groupID),
	}
}

func (c *Consumer) Subscribe(ctx context.Context, handler func([]byte) error) {
	go func() {
		for {
			m, err := c.r.ReadMessage(ctx)
			if err != nil {
				log.Println("Error reading message:", err)
				continue
			}

			if err := handler(m.Value); err != nil {
				log.Printf("Handler error, message: %s, err: %v\n", string(m.Value), err)
			}
		}
	}()
}

func (c *Consumer) Close() error {
	return c.r.Close()
}

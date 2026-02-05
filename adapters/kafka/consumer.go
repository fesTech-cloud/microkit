// adapters/kafka/consumer.go
package kafka

import (
	"context"
	"log"

	kafka "github.com/segmentio/kafka-go"
	"github.com/festech-cloud/microkit/internal/retry"
)

type ConsumerConfig struct {
	RetryConfig retry.Config
	EnableDLQ   bool
	DLQTopic    string
}

type Consumer struct {
	conn    *Connection
	topic   string
	groupID string
	r       *kafka.Reader
	config  ConsumerConfig
}

func NewConsumer(conn *Connection, topic, groupID string) *Consumer {
	return &Consumer{
		conn:    conn,
		topic:   topic,
		groupID: groupID,
		r:       conn.Reader(topic, groupID),
	}
}

func NewConsumerWithConfig(conn *Connection, topic, groupID string, config ConsumerConfig) *Consumer {
	return &Consumer{
		conn:    conn,
		topic:   topic,
		groupID: groupID,
		r:       conn.Reader(topic, groupID),
		config:  config,
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

			if c.config.RetryConfig.MaxAttempts > 0 {
				err = retry.Execute(ctx, c.config.RetryConfig, func() error {
					return handler(m.Value)
				})
			} else {
				err = handler(m.Value)
			}

			if err != nil {
				log.Printf("Handler error, message: %s, err: %v\n", string(m.Value), err)
				if c.config.EnableDLQ {
					c.sendToDLQ(ctx, m.Value)
				}
			}
		}
	}()
}

func (c *Consumer) sendToDLQ(ctx context.Context, message []byte) {
	producer := NewProducer(c.conn, c.config.DLQTopic)
	defer producer.Close()
	producer.Publish(ctx, nil, message)
}

func (c *Consumer) Close() error {
	return c.r.Close()
}
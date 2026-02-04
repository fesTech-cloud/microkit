package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	conn  *Connection
	topic string
	W     *kafka.Writer
}

func NewProducer(conn *Connection, topic string) *Producer {
	return &Producer{
		conn:  conn,
		topic: topic,
		W:     conn.Writer(topic),
	}
}

func (p *Producer) Publish(ctx context.Context, message []byte, key []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: message,
		Time:  time.Now(),
	}

	return p.W.WriteMessages(ctx, msg)
}

func (p *Producer) Close() error {
	return p.W.Close()
}

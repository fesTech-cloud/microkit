package kafka

import "github.com/segmentio/kafka-go"

type Connection struct {
	Brokers []string
}

func NewConnection(brokers []string) *Connection {
	return &Connection{
		Brokers: brokers,
	}
}

func (c *Connection) Writer(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(c.Brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}
}

func (c *Connection) Reader(topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.Brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,    // 1KB
		MaxBytes: 10e6, // 10MB
	})
}

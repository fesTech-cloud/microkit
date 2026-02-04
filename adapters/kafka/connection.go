package kafka

import (
	"context"
	"net"
	"strconv"

	sarama "github.com/IBM/sarama"
	kafka "github.com/segmentio/kafka-go"
)

type Connection struct {
	Brokers []string
	Config  *sarama.Config
}

func NewConnection(brokers []string) *Connection {
	cfg := sarama.NewConfig()

	// Required for producers
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll

	// Required for consumers
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	cfg.Version = sarama.V2_8_0_0
	return &Connection{
		Brokers: brokers,
		Config:  cfg,
	}
}

func (c *Connection) Writer(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(c.Brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           kafka.RequireAll,
		Async:                  false,
		AllowAutoTopicCreation: true,
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

func (c *Connection) CreateTopic(ctx context.Context, topic string) error {
	conn, err := kafka.Dial("tcp", c.Brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	return controllerConn.CreateTopics(topicConfigs...)
}

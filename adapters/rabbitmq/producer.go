package rabbitmq

import (
	"context"

	"github.com/festus/microkit/messaging"
	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	conn   *Connection
	ch     *amqp091.Channel
	config messaging.Config
}

func NewProducer(conn *Connection, cfg messaging.Config) (*Producer, error) {
	ch, err := conn.GetConnection().Channel()
	if err != nil {
		return nil, err
	}

	return &Producer{
		conn:   conn,
		ch:     ch,
		config: cfg,
	}, nil
}

func (p *Producer) Publish(ctx context.Context, topic string, msg messaging.Message) error {
	return p.ch.PublishWithContext(
		ctx,
		"amq.topic", // exchange
		topic,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        msg.Payload,
		},
	)
}

func (p *Producer) Close() error {
	if p.ch != nil {
		return p.ch.Close()
	}
	return nil
}

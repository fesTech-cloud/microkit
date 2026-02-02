package rabbitmq

import (
	"context"
	"log"

	"github.com/festech-cloud/microkit/messaging"
	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn   *Connection
	ch     *amqp091.Channel
	config messaging.Config
}

func NewConsumer(conn *Connection, cfg messaging.Config) (*Consumer, error) {
	ch, err := conn.GetConnection().Channel()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn:   conn,
		ch:     ch,
		config: cfg,
	}, nil
}

func (c *Consumer) Subscribe(ctx context.Context, topic string, handler messaging.HandlerFunc) error {
	queue, err := c.ch.QueueDeclare(
		topic,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	msgs, err := c.ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			msg := messaging.Message{Payload: d.Body, Headers: map[string]string{}}
			err := handler(ctx, msg)
			if err != nil {
				d.Nack(false, true)
				log.Println("Message handling failed, message requeued:", err)
			} else {
				d.Ack(false)
			}
		}
	}()

	return nil
}

func (c *Consumer) Close() error {
	if c.ch != nil {
		return c.ch.Close()
	}
	return nil
}

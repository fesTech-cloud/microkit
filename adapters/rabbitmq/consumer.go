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
	// ---- DLQ setup ----
	dlxName := topic + ".dlx"
	dlqName := topic + ".dlq"

	// Dead-letter exchange
	if err := c.ch.ExchangeDeclare(
		dlxName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	// Dead-letter queue
	if _, err := c.ch.QueueDeclare(
		dlqName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	// Bind DLQ to DLX
	if err := c.ch.QueueBind(
		dlqName,
		topic,
		dlxName,
		false,
		nil,
	); err != nil {
		return err
	}

	// Main queue with DLQ configuration
	queueArgs := amqp091.Table{
		"x-dead-letter-exchange":    dlxName,
		"x-dead-letter-routing-key": topic,
	}

	queue, err := c.ch.QueueDeclare(
		topic,
		true,
		false,
		false,
		false,
		queueArgs,
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
				d.Nack(false, false) // goes to DLQ
				log.Printf("Message processing failed, sent to DLQ: %v", err)
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

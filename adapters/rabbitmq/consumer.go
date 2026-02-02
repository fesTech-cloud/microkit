package rabbitmq

import (
	"context"
	"log"

	"github.com/festech-cloud/microkit/messaging"
	"github.com/rabbitmq/amqp091-go"
)

var (
	maxRetries = 3
	retryDelay = 5000 // milliseconds
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

// Public API

func (c *Consumer) Subscribe(
	ctx context.Context,
	topic string,
	handler messaging.HandlerFunc,
) error {
	dlx, err := c.setupDLX(topic)
	if err != nil {
		return err
	}

	if err := c.setupDLQ(topic, dlx); err != nil {
		return err
	}

	if err := c.setupRetryQueue(topic); err != nil {
		return err
	}

	queue, err := c.setupMainQueue(topic, dlx)
	if err != nil {
		return err
	}

	msgs, err := c.consume(queue)
	if err != nil {
		return err
	}

	c.handleMessages(ctx, msgs, topic, handler)
	return nil
}

func (c *Consumer) Close() error {
	if c.ch != nil {
		return c.ch.Close()
	}
	return nil
}

// Queue setup helpers

func (c *Consumer) setupDLX(topic string) (string, error) {
	name := dlxName(topic)

	return name, c.ch.ExchangeDeclare(
		name,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (c *Consumer) setupDLQ(topic, dlx string) error {
	_, err := c.ch.QueueDeclare(
		dlqName(topic),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return c.ch.QueueBind(
		dlqName(topic),
		topic,
		dlx,
		false,
		nil,
	)
}

func (c *Consumer) setupRetryQueue(topic string) error {
	_, err := c.ch.QueueDeclare(
		retryName(topic),
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-message-ttl":             int32(retryDelay),
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": topic,
		},
	)
	return err
}

func (c *Consumer) setupMainQueue(topic, dlx string) (amqp091.Queue, error) {
	queue, err := c.ch.QueueDeclare(
		topic,
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-dead-letter-exchange":    dlx,
			"x-dead-letter-routing-key": topic,
		},
	)

	if err != nil {
		return amqp091.Queue{}, err
	}

	err = c.ch.QueueBind(
		queue.Name,
		topic,
		"amq.topic",
		false,
		nil,
	)
	if err != nil {
		return amqp091.Queue{}, err
	}

	return queue, nil
}

func (c *Consumer) consume(queue amqp091.Queue) (<-chan amqp091.Delivery, error) {
	return c.ch.Consume(
		queue.Name,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
}

// Message handling

func (c *Consumer) handleMessages(
	ctx context.Context,
	msgs <-chan amqp091.Delivery,
	topic string,
	handler messaging.HandlerFunc,
) {
	go func() {
		for d := range msgs {
			retries := getRetryCount(d.Headers)

			msg := messaging.Message{
				Payload: d.Body,
				Headers: map[string]string{},
			}

			if err := handler(ctx, msg); err != nil {
				if retries >= maxRetries {
					log.Println("Max retries exceeded, sending to DLQ:", err)
					d.Nack(false, false)
					continue
				}

				headers := d.Headers
				if headers == nil {
					headers = amqp091.Table{}
				}
				headers["x-retry-count"] = retries + 1

				_ = c.ch.Publish(
					"",
					retryName(topic),
					false,
					false,
					amqp091.Publishing{
						Headers: headers,
						Body:    d.Body,
					},
				)

				d.Ack(false)
			} else {
				d.Ack(false)
			}
		}
	}()
}

// Helpers

func getRetryCount(headers amqp091.Table) int {
	if headers == nil {
		return 0
	}

	if val, ok := headers["x-retry-count"]; ok {
		switch v := val.(type) {
		case int32:
			return int(v)
		case int:
			return v
		}
	}
	return 0
}

func dlxName(topic string) string {
	return topic + ".dlx"
}

func dlqName(topic string) string {
	return topic + ".dlq"
}

func retryName(topic string) string {
	return topic + ".retry"
}

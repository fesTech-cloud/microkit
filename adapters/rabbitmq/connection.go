package rabbitmq

import (
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	URL  string
	conn *amqp091.Connection
}

func NewConnection(url string) (*Connection, error) {
	c := &Connection{URL: url}

	err := c.connect()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Connection) connect() error {
	var err error

	// Retry connection up to 3 times
	for range 3 {
		c.conn, err = amqp091.Dial(c.URL)
		if err == nil {
			return nil
		}
		log.Println("RabbitMQ connection failed, retrying in 1s...")
		time.Sleep(1 * time.Second)
	}
	return err
}

func (c *Connection) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Connection) GetConnection() *amqp091.Connection {
	return c.conn
}

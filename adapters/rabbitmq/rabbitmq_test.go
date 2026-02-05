package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/festus/microkit/messaging"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// -------------------------
// Globals for tests
// -------------------------
var testConn *Connection
var ctx = context.Background()

// -------------------------
// Start RabbitMQ container
// -------------------------
func setupRabbitMQ(t *testing.T) string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        "rabbitmq:3.11-alpine", // lighter image for tests
		ExposedPorts: []string{"5672/tcp"},
		WaitingFor:   wait.ForListeningPort("5672/tcp").WithStartupTimeout(2 * time.Minute),
	}

	rmqC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start RabbitMQ container: %v", err)
	}

	// Ensure container cleanup after test
	t.Cleanup(func() {
		_ = rmqC.Terminate(ctx)
	})

	port, err := rmqC.MappedPort(ctx, "5672")
	if err != nil {
		t.Fatalf("Failed to get mapped port: %v", err)
	}

	url := fmt.Sprintf("amqp://guest:guest@localhost:%s/", port.Port())
	return url
}

// -------------------------
// Helper to create connection with retry
// -------------------------
func initConnection(t *testing.T) *Connection {
	url := setupRabbitMQ(t)

	var conn *Connection
	var err error
	for i := 0; i < 5; i++ {
		conn, err = NewConnection(url)
		if err == nil {
			break
		}
		log.Println("Retrying connection to RabbitMQ...")
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ after retries: %v", err)
	}
	return conn
}

// -------------------------
// Test: Connection
// -------------------------
func TestConnection(t *testing.T) {
	conn := initConnection(t)
	defer conn.Close()

	if conn.conn.IsClosed() {
		t.Fatal("Connection should be open")
	}
	t.Log("Connection test passed")
}

// -------------------------
// Test: Producer + Consumer
// -------------------------
func TestProducerConsumer(t *testing.T) {
	conn := initConnection(t)
	defer conn.Close()

	topic := "test-topic"

	// Create Producer
	prodCfg := messaging.Config{}
	producer, err := NewProducer(conn, prodCfg)
	if err != nil {
		t.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Create Consumer
	consCfg := messaging.Config{}
	consumer, err := NewConsumer(conn, consCfg)
	if err != nil {
		t.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Channel to receive message
	received := make(chan []byte, 1)

	// Start consumer
	err = consumer.Subscribe(ctx, topic, func(ctx context.Context, msg messaging.Message) error {
		received <- msg.Payload
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to subscribe consumer: %v", err)
	}

	// Publish a message
	message := messaging.Message{Payload: []byte("hello world")}
	err = producer.Publish(ctx, topic, message)
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	// Wait for message to be received
	select {
	case msg := <-received:
		if string(msg) != "hello world" {
			t.Fatalf("Expected 'hello world', got '%s'", string(msg))
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for message")
	}

	t.Log("Producer + Consumer test passed")
}

// -------------------------
// Test: Retry / DLQ simulation
// -------------------------
func TestRetryDLQ(t *testing.T) {
	conn := initConnection(t)
	defer conn.Close()

	topic := "retry-topic"

	prodCfg := messaging.Config{}
	producer, err := NewProducer(conn, prodCfg)
	if err != nil {
		t.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	consCfg := messaging.Config{}
	consumer, err := NewConsumer(conn, consCfg)
	if err != nil {
		t.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Reduce retry delay for testing
	retryDelay = 500 // 500ms instead of 5000ms

	// Channel to track retries
	retryCount := 0
	maxRetries := 2

	err = consumer.Subscribe(ctx, topic, func(ctx context.Context, msg messaging.Message) error {
		retryCount++
		if retryCount <= maxRetries {
			return fmt.Errorf("simulate failure %d", retryCount)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to subscribe consumer: %v", err)
	}

	// Publish a message
	message := messaging.Message{Payload: []byte("retry test")}
	err = producer.Publish(ctx, topic, message)
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	// Poll until retryCount reaches expected value or timeout
	timeout := time.After(5 * time.Second)
	tick := time.Tick(100 * time.Millisecond)

	for {
		select {
		case <-timeout:
			t.Fatalf("Expected %d retries, got %d", maxRetries, retryCount)
		case <-tick:
			if retryCount >= maxRetries+1 {
				t.Logf("Retry/DLQ simulation passed with %d attempts", retryCount)
				return
			}
		}
	}
}

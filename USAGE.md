# Installation

```bash
go get github.com/festus/microkit
```

# Quick Start

## HTTP Client

```go
package main

import (
    "context"
    "time"
    
    "github.com/festus/microkit/adapters/http"
    "github.com/festus/microkit/network"
)

func main() {
    client := http.NewClient(10 * time.Second)
    defer client.Close()
    
    resp, err := client.Get(context.Background(), "https://api.example.com",
        network.WithHeader("Authorization", "Bearer token"),
        network.WithRetry(3, 100*time.Millisecond, 2*time.Second, 2.0),
    )
    // Handle response...
}
```

## Kafka Producer

```go
package main

import (
    "context"
    
    "github.com/festus/microkit/adapters/kafka"
    "github.com/festus/microkit/messaging"
)

func main() {
    conn := kafka.NewConnection([]string{"localhost:9092"})
    producer := kafka.NewProducer(conn, "my-topic")
    defer producer.Close()
    
    producer.Publish(context.Background(), []byte("key"), []byte("message"))
}
```

## RabbitMQ Consumer

```go
package main

import (
    "context"
    
    "github.com/festus/microkit/adapters/rabbitmq"
    "github.com/festus/microkit/messaging"
)

func main() {
    producer, _ := rabbitmq.NewProducer(
        rabbitmq.WithURL("amqp://localhost"),
    )
    defer producer.Close()
    
    producer.Publish(context.Background(), "queue.name", messaging.Message{
        Payload: []byte(`{"data": "value"}`),
    })
}
```
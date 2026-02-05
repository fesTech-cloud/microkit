# microkit

[![Go Reference](https://pkg.go.dev/badge/github.com/festus/microkit.svg)](https://pkg.go.dev/github.com/festus/microkit)
[![Go Report Card](https://goreportcard.com/badge/github.com/festus/microkit)](https://goreportcard.com/report/github.com/festus/microkit)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

An opinionated Go toolkit that simplifies building microservices by providing clean abstractions over common infrastructure tools like messaging, networking, and service communication.

---

## Why microkit?

Building microservices in Go often means rewriting the same infrastructure code:

- Messaging setup
- Network calls
- Retries and timeouts
- Graceful shutdown
- Error handling

**microkit** standardizes these patterns with clear, minimal abstractions so developers can focus on business logic instead of boilerplate.

This is not a framework.
It’s a toolkit you can adopt incrementally.

---

## Getting Started

### Installation

```bash
go get github.com/festus/microkit
```

### Quick Start

#### HTTP Client with Retry

```go
package main

import (
    "context"
    "fmt"
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
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Response: %d\n", resp.StatusCode)
}
```

#### Kafka Consumer with Retry/DLQ

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/festus/microkit/adapters/kafka"
    "github.com/festus/microkit/internal/retry"
)

func main() {
    conn := kafka.NewConnection([]string{"localhost:9092"})
    
    config := kafka.ConsumerConfig{
        RetryConfig: retry.Config{
            MaxAttempts:  3,
            InitialDelay: 100 * time.Millisecond,
            MaxDelay:     2 * time.Second,
            Multiplier:   2.0,
        },
        EnableDLQ: true,
        DLQTopic:  "orders-dlq",
    }
    
    consumer := kafka.NewConsumerWithConfig(conn, "orders", "order-service", config)
    defer consumer.Close()
    
    consumer.Subscribe(context.Background(), func(msg []byte) error {
        fmt.Printf("Processing: %s\n", string(msg))
        // Your business logic here
        return nil
    })
    
    select {} // Keep running
}
```

---

## Current Focus

- Messaging abstractions
- Network abstractions (HTTP & gRPC)
- RabbitMQ adapter (v0.1)
- Kafka adapter (v0.1)
- HTTP adapter (v0.1)
- gRPC adapter (v0.1)

---

## Examples

### RabbitMQ

```go
producer, _ := rabbitmq.NewProducer(
    rabbitmq.WithURL("amqp://localhost"),
)

producer.Publish(ctx, "orders.created", messaging.Message{
    Payload: []byte(`{"id":"123"}`),
})
```

### Kafka

```go
conn := kafka.NewConnection([]string{"localhost:9092"})
producer := kafka.NewProducer(conn, "orders")

producer.Publish(ctx, []byte("key"), []byte(`{"id":"123"}`))
```

### HTTP Client with Retry

```go
client := http.NewClient(10 * time.Second)
defer client.Close()

resp, _ := client.Get(ctx, "https://api.example.com/users",
    network.WithHeader("Authorization", "Bearer token"),
    network.WithRetry(3, 100*time.Millisecond, 2*time.Second, 2.0),
)
```

### gRPC Client with Retry

```go
client, _ := grpc.NewClient("localhost:50051", 10*time.Second)
defer client.Close()

req := &HelloRequest{Name: "microkit"}
resp := &HelloResponse{}

client.Call(ctx, "/hello.HelloService/SayHello", req, resp,
    network.WithRetry(3, 100*time.Millisecond, 2*time.Second, 2.0),
)
```

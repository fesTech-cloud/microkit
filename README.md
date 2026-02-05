# microkit

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

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
- RabbitMQ adapter (v0.1)
- Kafka adapter (v0.1)

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

package messaging

import "context"

// Producer defines the interface for publishing messages.
type Producer interface {
	Publish(ctx context.Context, topic string, msg Message) error
	Close() error
}

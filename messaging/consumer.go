package messaging

import "context"

// Consumer defines the interface for consuming messages.
type Consumer interface {
	Subscribe(ctx context.Context, topic string, handler HandlerFunc) error
	Close() error
}

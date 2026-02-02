package messaging

import "context"

// HandlerFunc defines the function signature for consuming messages.
type HandlerFunc func(ctx context.Context, msg Message) error

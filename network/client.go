package network

import "context"

// Client defines the interface for making network calls (HTTP/gRPC).
type Client interface {
	// HTTP methods
	Get(ctx context.Context, url string, opts ...Option) (*Response, error)
	Post(ctx context.Context, url string, body []byte, opts ...Option) (*Response, error)
	Put(ctx context.Context, url string, body []byte, opts ...Option) (*Response, error)
	Delete(ctx context.Context, url string, opts ...Option) (*Response, error)
	
	// gRPC method
	Call(ctx context.Context, method string, req interface{}, resp interface{}, opts ...Option) error
	
	Close() error
}
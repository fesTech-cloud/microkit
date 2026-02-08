package network

import "errors"

// Response represents a network response.
type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

// Request represents a network request.
type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
}

// Common errors
var (
	ErrUnsupportedOperation = errors.New("operation not supported by this client")
)
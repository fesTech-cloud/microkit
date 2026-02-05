package network

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
package messaging

// Message represents a single message in the system.
type Message struct {
	ID        string
	Payload   []byte
	Headers   map[string]string
	Timestamp int64
}

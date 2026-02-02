package messaging

import "time"

// Config holds library-wide options for producers and consumers.
type Config struct {
	RetryCount int
	RetryDelay time.Duration
	Timeout    time.Duration
}

// DefaultConfig returns a reasonable default configuration
func DefaultConfig() Config {
	return Config{
		RetryCount: 3,
		RetryDelay: 500 * time.Millisecond,
		Timeout:    5 * time.Second,
	}
}

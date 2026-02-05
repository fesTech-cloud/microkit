package network

import (
	"time"

	"github.com/festech-cloud/microkit/internal/retry"
)

// Option configures network requests.
type Option func(*Config)

// Config holds configuration for network requests.
type Config struct {
	Headers     map[string]string
	Timeout     time.Duration
	RetryConfig *retry.Config
}

// WithHeader adds a header to the request.
func WithHeader(key, value string) Option {
	return func(c *Config) {
		if c.Headers == nil {
			c.Headers = make(map[string]string)
		}
		c.Headers[key] = value
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithRetry enables retry logic for network calls.
func WithRetry(maxAttempts int, initialDelay, maxDelay time.Duration, multiplier float64) Option {
	return func(c *Config) {
		c.RetryConfig = &retry.Config{
			MaxAttempts:  maxAttempts,
			InitialDelay: initialDelay,
			MaxDelay:     maxDelay,
			Multiplier:   multiplier,
		}
	}
}
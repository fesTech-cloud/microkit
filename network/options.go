package network

import (
	"time"

	"github.com/festus/microkit/internal/retry"
)

// Option configures network requests.
type Option func(*Config)

// Config holds configuration for network requests.
type Config struct {
	Headers          map[string]string
	Timeout          time.Duration
	RetryConfig      *retry.Config
	RetryStatusCodes []int
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
// Set jitter to true to add randomization and prevent thundering herd.
func WithRetry(maxAttempts int, initialDelay, maxDelay time.Duration, multiplier float64) Option {
	return func(c *Config) {
		c.RetryConfig = &retry.Config{
			MaxAttempts:  maxAttempts,
			InitialDelay: initialDelay,
			MaxDelay:     maxDelay,
			Multiplier:   multiplier,
			Jitter:       true,
		}
	}
}

// WithRetryOnStatus enables retry based on HTTP status codes.
func WithRetryOnStatus(maxAttempts int, initialDelay, maxDelay time.Duration, multiplier float64, statusCodes ...int) Option {
	return func(c *Config) {
		c.RetryConfig = &retry.Config{
			MaxAttempts:  maxAttempts,
			InitialDelay: initialDelay,
			MaxDelay:     maxDelay,
			Multiplier:   multiplier,
			Jitter:       true,
		}
		c.RetryStatusCodes = statusCodes
	}
}
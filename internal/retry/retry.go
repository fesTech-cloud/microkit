package retry

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

var ErrRetryableStatus = errors.New("retryable status code")

type Config struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
	Jitter       bool
}

func Execute(ctx context.Context, config Config, fn func() error) error {
	var err error
	delay := config.InitialDelay

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}

		if attempt < config.MaxAttempts-1 {
			actualDelay := delay
			if config.Jitter {
				actualDelay = time.Duration(float64(delay) * (0.5 + rand.Float64()*0.5))
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(actualDelay):
			}

			delay = time.Duration(float64(delay) * config.Multiplier)
			if delay > config.MaxDelay {
				delay = config.MaxDelay
			}
		}
	}

	return err
}

package withretry

import (
	"context"
	"time"

	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
)

type RetryConfig struct {
	MaxAttempts       uint
	InitialBackoff    time.Duration
	BackoffMultiplier float64
}

var DefaultConfig = RetryConfig{
	MaxAttempts:       5,
	InitialBackoff:    time.Millisecond * 100,
	BackoffMultiplier: 2,
}

func WithRetry(
	ctx context.Context,
	config RetryConfig,
	isTransient func(error) bool,
	operation func(context.Context) error,
) error {
	backoff := config.InitialBackoff

	for attempt := uint(0); attempt < config.MaxAttempts; attempt++ {
		err := operation(ctx)
		if err == nil {
			return nil
		}

		if !isTransient(err) {
			return err
		}

		if attempt == config.MaxAttempts-1 {
			break
		}

		select {
		case <-ctx.Done():
			return apierrors.RequestTimeout("request timeout", ctx.Err().Error())
		case <-time.After(backoff):
			backoff = time.Duration(float64(backoff) * config.BackoffMultiplier)
		}
	}

	return apierrors.RequestTimeout("max retry attempts reached")
}

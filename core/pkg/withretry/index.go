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

func WithRetry[T any](
	ctx context.Context,
	config RetryConfig,
	isTransient func(error) bool,
	operation func(context.Context) (T, error),
) (T, error) {
	var zero T
	backoff := config.InitialBackoff

	for attempt := uint(0); attempt < config.MaxAttempts; attempt++ {
		result, err := operation(ctx)
		if err == nil {
			return result, nil
		}

		if !isTransient(err) {
			return zero, err
		}

		if attempt == config.MaxAttempts-1 {
			break
		}

		select {
		case <-ctx.Done():
			return zero, apierrors.RequestTimeout("request timeout", ctx.Err().Error())
		case <-time.After(backoff):
			backoff = time.Duration(float64(backoff) * config.BackoffMultiplier)
		}
	}

	return zero, apierrors.RequestTimeout("max retry attempts reached")
}

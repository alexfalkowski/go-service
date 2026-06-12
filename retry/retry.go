package retry

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/time"
	retry "github.com/sethvargo/go-retry"
)

// Backoff aliases the upstream retry backoff interface.
type Backoff = retry.Backoff

// RetryFunc aliases the upstream retry function type.
type RetryFunc = retry.RetryFunc

// RetryFuncValue aliases the upstream retry value function type.
type RetryFuncValue[T any] = retry.RetryFuncValue[T]

// NewConstant creates a constant backoff.
func NewConstant(duration time.Duration) Backoff {
	return retry.NewConstant(duration.Duration())
}

// WithJitterPercent wraps next and adds +/- percent jitter to each backoff duration.
func WithJitterPercent(percent uint64, next Backoff) Backoff {
	return retry.WithJitterPercent(percent, next)
}

// RetryableError marks err as retryable.
func RetryableError(err error) error {
	return retry.RetryableError(err)
}

// WithMaxRetries wraps next so it retries at most attempts times after the initial attempt.
func WithMaxRetries(attempts uint64, next Backoff) Backoff {
	return retry.WithMaxRetries(attempts, next)
}

// Do wraps f with b and retries retryable errors.
func Do(ctx context.Context, b Backoff, f RetryFunc) error {
	return retry.Do(ctx, b, f)
}

// DoValue wraps f with b and retries retryable errors, returning the produced value.
func DoValue[T any](ctx context.Context, b Backoff, f RetryFuncValue[T]) (T, error) {
	return retry.DoValue(ctx, b, f)
}

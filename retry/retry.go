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

// NewBackoff creates a backoff for the given strategy using base as the starting duration.
//
// Supported strategies are "constant" (the default), "exponential", and "fibonacci".
// An empty or unrecognized strategy falls back to a constant backoff. It forwards to the
// upstream retry package and panics when base is less than or equal to zero.
func NewBackoff(strategy string, base time.Duration) Backoff {
	switch strategy {
	case "exponential":
		return retry.NewExponential(base.Duration())
	case "fibonacci":
		return retry.NewFibonacci(base.Duration())
	default:
		return retry.NewConstant(base.Duration())
	}
}

// WithJitterPercent wraps next and adds +/- percent jitter to each backoff duration.
//
// The returned backoff delegates stop behavior to next.
func WithJitterPercent(percent uint64, next Backoff) Backoff {
	return retry.WithJitterPercent(percent, next)
}

// RetryableError marks err as retryable.
//
// It returns nil when err is nil. When retries are exhausted, [Do] and [DoValue]
// return the original wrapped cause rather than the retry marker.
func RetryableError(err error) error {
	return retry.RetryableError(err)
}

// WithMaxRetries wraps next so it allows at most attempts retry waits after the initial call.
//
// Pass zero to stop after the first retryable error.
func WithMaxRetries(attempts uint64, next Backoff) Backoff {
	return retry.WithMaxRetries(attempts, next)
}

// WithCappedDuration wraps next so each returned backoff duration is capped at maxDuration.
//
// This bounds only the per-attempt duration, not the total backoff time; combine it with
// [WithMaxRetries] or [Do]'s context to bound overall retry time.
func WithCappedDuration(maxDuration time.Duration, next Backoff) Backoff {
	return retry.WithCappedDuration(maxDuration.Duration(), next)
}

// Do wraps f with b and retries errors marked by [RetryableError].
//
// It forwards ctx to each call of f. If ctx is canceled before or during a
// backoff wait, Do returns ctx.Err().
func Do(ctx context.Context, b Backoff, f RetryFunc) error {
	return retry.Do(ctx, b, f)
}

// DoValue wraps f with b and retries errors marked by [RetryableError], returning the produced value.
//
// It forwards ctx to each call of f. If ctx is canceled before or during a
// backoff wait, DoValue returns the zero value of T with ctx.Err().
func DoValue[T any](ctx context.Context, b Backoff, f RetryFuncValue[T]) (T, error) {
	return retry.DoValue(ctx, b, f)
}

package breaker

import (
	"github.com/alexfalkowski/go-service/v2/time"
	breaker "github.com/sony/gobreaker"
)

// NewCircuitBreaker constructs a circuit breaker using the provided settings.
func NewCircuitBreaker(st Settings) *CircuitBreaker {
	return breaker.NewCircuitBreaker(st)
}

// CircuitBreaker is an alias for breaker.CircuitBreaker.
type CircuitBreaker = breaker.CircuitBreaker

// Counts is an alias for breaker.Counts.
type Counts = breaker.Counts

// Settings is an alias for breaker.Settings.
type Settings = breaker.Settings

// DefaultSettings provides a conservative default breaker configuration.
//
// It allows up to 3 requests when half-open, uses a 30s metrics window, a 10s open timeout,
// and trips after 5 consecutive failures.
var DefaultSettings = Settings{
	MaxRequests: 3,
	Interval:    30 * time.Second,
	Timeout:     10 * time.Second,
	ReadyToTrip: func(counts breaker.Counts) bool {
		return counts.ConsecutiveFailures >= 5
	},
}

// ErrOpenState is an alias for breaker.ErrOpenState.
var ErrOpenState = breaker.ErrOpenState

// ErrTooManyRequests is an alias for breaker.ErrTooManyRequests.
var ErrTooManyRequests = breaker.ErrTooManyRequests

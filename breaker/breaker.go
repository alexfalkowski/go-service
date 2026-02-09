package breaker

import (
	"github.com/alexfalkowski/go-service/v2/time"
	breaker "github.com/sony/gobreaker"
)

// NewCircuitBreaker is an alias for the breaker.NewCircuitBreaker.
func NewCircuitBreaker(st Settings) *CircuitBreaker {
	return breaker.NewCircuitBreaker(st)
}

// CircuitBreaker is an alias for the breaker.CircuitBreaker.
type CircuitBreaker = breaker.CircuitBreaker

// Counts is an alias for the breaker.Counts.
type Counts = breaker.Counts

// Settings is an alias for the breaker.Settings.
type Settings = breaker.Settings

// DefaultSettings is the default settings for the breaker.
var DefaultSettings = Settings{
	MaxRequests: 3,
	Interval:    30 * time.Second,
	Timeout:     10 * time.Second,
	ReadyToTrip: func(counts breaker.Counts) bool {
		return counts.ConsecutiveFailures >= 5
	},
}

// ErrOpenState is an alias for the breaker.ErrOpenState.
var ErrOpenState = breaker.ErrOpenState

// ErrTooManyRequests is an alias for the breaker.ErrTooManyRequests.
var ErrTooManyRequests = breaker.ErrTooManyRequests

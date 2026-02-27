package breaker

import (
	"github.com/alexfalkowski/go-service/v2/time"
	breaker "github.com/sony/gobreaker"
)

// NewCircuitBreaker constructs a new circuit breaker using the provided Settings.
//
// This is a thin wrapper around github.com/sony/gobreaker.NewCircuitBreaker that exists so
// go-service can:
//   - re-export gobreaker types behind a stable package path, and
//   - centralize shared defaults via DefaultSettings.
//
// Callers typically customize Settings (e.g. Name, ReadyToTrip, IsSuccessful) and then pass it here.
func NewCircuitBreaker(st Settings) *CircuitBreaker {
	return breaker.NewCircuitBreaker(st)
}

// CircuitBreaker is an alias for github.com/sony/gobreaker.CircuitBreaker.
//
// It is a generic circuit breaker implementation used by go-service transports. Prefer importing
// this package's alias when interacting with breakers created by NewCircuitBreaker.
type CircuitBreaker = breaker.CircuitBreaker

// Counts is an alias for github.com/sony/gobreaker.Counts.
//
// Counts is used by Settings.ReadyToTrip to decide whether the breaker should open.
type Counts = breaker.Counts

// Settings is an alias for github.com/sony/gobreaker.Settings.
//
// Settings controls breaker behavior (half-open behavior, rolling interval, open timeout, and
// success/failure classification via IsSuccessful).
type Settings = breaker.Settings

// DefaultSettings provides a conservative default circuit breaker configuration.
//
// Defaults:
//   - MaxRequests: 3 (allowed while half-open)
//   - Interval: 30s (rolling window for internal Counts)
//   - Timeout: 10s (time breaker stays open before transitioning to half-open)
//   - ReadyToTrip: open after 5 consecutive failures
//
// Transport integrations typically copy DefaultSettings and then set:
//   - Settings.Name to a stable per-upstream/per-method key
//   - Settings.IsSuccessful to classify which errors should count as failures
var DefaultSettings = Settings{
	MaxRequests: 3,
	Interval:    30 * time.Second,
	Timeout:     10 * time.Second,
	ReadyToTrip: func(counts breaker.Counts) bool {
		return counts.ConsecutiveFailures >= 5
	},
}

// ErrOpenState is an alias for github.com/sony/gobreaker.ErrOpenState.
//
// It is returned by CircuitBreaker.Execute when the breaker is open.
var ErrOpenState = breaker.ErrOpenState

// ErrTooManyRequests is an alias for github.com/sony/gobreaker.ErrTooManyRequests.
//
// It is returned by CircuitBreaker.Execute when the breaker is half-open and MaxRequests
// would be exceeded.
var ErrTooManyRequests = breaker.ErrTooManyRequests

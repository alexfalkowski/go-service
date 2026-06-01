package breaker_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/breaker"
	"github.com/stretchr/testify/require"
)

func TestDefaultSettings(t *testing.T) {
	settings := breaker.DefaultSettings

	require.Equal(t, uint32(3), settings.MaxRequests)
	require.Equal(t, (30 * time.Second).Duration(), settings.Interval)
	require.Equal(t, (10 * time.Second).Duration(), settings.Timeout)
	require.False(t, settings.ReadyToTrip(breaker.Counts{ConsecutiveFailures: 4}))
	require.True(t, settings.ReadyToTrip(breaker.Counts{ConsecutiveFailures: 5}))
}

func TestNewCircuitBreaker(t *testing.T) {
	errFailed := errors.New("failed")
	cb := breaker.NewCircuitBreaker(breaker.Settings{
		ReadyToTrip: func(counts breaker.Counts) bool {
			return counts.ConsecutiveFailures >= 1
		},
	})

	_, err := cb.Execute(func() (any, error) {
		return nil, errFailed
	})
	require.ErrorIs(t, err, errFailed)

	_, err = cb.Execute(func() (any, error) {
		return "ok", nil
	})
	require.ErrorIs(t, err, breaker.ErrOpenState)
}

func TestStateChangeUsesStableTypes(t *testing.T) {
	errFailed := errors.New("failed")
	var from, to breaker.State
	cb := breaker.NewCircuitBreaker(breaker.Settings{
		ReadyToTrip: func(counts breaker.Counts) bool {
			return counts.ConsecutiveFailures >= 1
		},
		OnStateChange: func(_ string, f breaker.State, t breaker.State) {
			from = f
			to = t
		},
	})

	_, err := cb.Execute(func() (any, error) {
		return nil, errFailed
	})

	require.ErrorIs(t, err, errFailed)
	require.Equal(t, breaker.StateClosed, from, "state change should start from closed")
	require.Equal(t, breaker.StateOpen, to, "state change should transition to open")
	require.Equal(t, "open", breaker.StateOpen.String(), "state alias should preserve upstream methods")
}

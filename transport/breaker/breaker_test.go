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

	tests := []struct {
		name   string
		counts breaker.Counts
		want   bool
	}{
		{
			name:   "below consecutive threshold",
			counts: breaker.Counts{ConsecutiveFailures: 4},
			want:   false,
		},
		{
			name: "total failures below consecutive threshold",
			counts: breaker.Counts{
				Requests:            6,
				TotalFailures:       5,
				ConsecutiveFailures: 4,
			},
			want: false,
		},
		{
			name:   "consecutive threshold",
			counts: breaker.Counts{ConsecutiveFailures: 5},
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, settings.ReadyToTrip(tt.counts))
		})
	}
}

func TestConfig(t *testing.T) {
	require.False(t, (*breaker.Config)(nil).IsEnabled())
	require.True(t, (&breaker.Config{}).IsEnabled())
}

func TestConfigSettings(t *testing.T) {
	settings := (&breaker.Config{
		MaxRequests:         2,
		Interval:            15 * time.Second,
		Timeout:             5 * time.Second,
		ConsecutiveFailures: 4,
	}).Settings()

	require.Equal(t, uint32(2), settings.MaxRequests)
	require.Equal(t, (15 * time.Second).Duration(), settings.Interval)
	require.Equal(t, (5 * time.Second).Duration(), settings.Timeout)
	require.False(t, settings.ReadyToTrip(breaker.Counts{ConsecutiveFailures: 3}))
	require.True(t, settings.ReadyToTrip(breaker.Counts{ConsecutiveFailures: 4}))
}

func TestConfigSettingsDefaults(t *testing.T) {
	settings := (*breaker.Config)(nil).Settings()

	require.Equal(t, breaker.DefaultSettings.MaxRequests, settings.MaxRequests)
	require.Equal(t, breaker.DefaultSettings.Interval, settings.Interval)
	require.Equal(t, breaker.DefaultSettings.Timeout, settings.Timeout)
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

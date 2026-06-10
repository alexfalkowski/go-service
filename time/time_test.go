package time_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestNewTimer(t *testing.T) {
	timer := time.NewTimer(time.Nanosecond)
	require.NotNil(t, timer)
	defer timer.Stop()

	select {
	case tm := <-timer.C:
		require.False(t, tm.IsZero())
	case <-time.After(time.Second):
		require.FailNow(t, "timed out waiting for timer")
	}
}

func TestUntil(t *testing.T) {
	future := time.Now().Add(time.Hour.Duration())

	duration := time.Until(future)

	require.Positive(t, duration)
	require.LessOrEqual(t, duration, time.Hour)
}

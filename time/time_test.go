package time_test

import (
	"testing"
	"testing/synctest"

	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestNewTimer(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		timer := time.NewTimer(time.Nanosecond)
		require.NotNil(t, timer)
		t.Cleanup(func() {
			timer.Stop()
		})

		tm := <-timer.C
		require.False(t, tm.IsZero())
	})
}

func TestUntil(t *testing.T) {
	future := time.Now().Add(time.Hour.Duration())

	duration := time.Until(future)

	require.Positive(t, duration)
	require.LessOrEqual(t, duration, time.Hour)
}

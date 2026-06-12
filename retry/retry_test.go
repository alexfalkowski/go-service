package retry_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/retry"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestWithJitterPercentBoundsBackoff(t *testing.T) {
	backoff := retry.WithJitterPercent(20, retry.NewConstant(100*time.Millisecond))

	durations := map[time.Duration]struct{}{}
	for range 100 {
		next, stop := backoff.Next()
		duration := time.Duration(next)

		require.False(t, stop, "jittered constant backoff should not stop")
		require.GreaterOrEqual(t, duration, 80*time.Millisecond, "jittered backoff should stay above -20%%")
		require.LessOrEqual(t, duration, 120*time.Millisecond, "jittered backoff should stay below +20%%")

		durations[duration] = struct{}{}
	}

	require.Greater(t, len(durations), 1, "jittered backoff should decorrelate repeated delays")
}

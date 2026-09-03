package retry_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/retry"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestNewBackoffProducesStrategyDelays(t *testing.T) {
	base := 100 * time.Millisecond

	tests := []struct {
		name     string
		strategy string
		want     []time.Duration
	}{
		{name: "empty defaults to constant", strategy: "", want: []time.Duration{base, base, base}},
		{name: "constant", strategy: "constant", want: []time.Duration{base, base, base}},
		{name: "unknown defaults to constant", strategy: "bogus", want: []time.Duration{base, base, base}},
		{name: "exponential", strategy: "exponential", want: []time.Duration{base, 2 * base, 4 * base}},
		{name: "fibonacci", strategy: "fibonacci", want: []time.Duration{base, 2 * base, 3 * base}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backoff := retry.NewBackoff(tt.strategy, base)

			for _, want := range tt.want {
				next, stop := backoff.Next()

				require.False(t, stop, "strategy backoff should not stop before max retries")
				require.Equal(t, want, time.Duration(next))
			}
		})
	}
}

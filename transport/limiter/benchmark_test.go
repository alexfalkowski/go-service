package limiter_test

import (
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/limiter"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

// BenchmarkLimiterManyDistinctKeys measures limiter behavior after its independent-key cap is reached.
func BenchmarkLimiterManyDistinctKeys(b *testing.B) {
	lc := fxtest.NewLifecycle(b)
	lim, err := limiter.NewLimiter(lc, limiter.KeyMap{"user-agent": meta.UserAgent}, &limiter.Config{
		Kind:     "user-agent",
		Tokens:   1,
		Interval: time.Second,
		MaxKeys:  limiter.DefaultMaxKeys,
	})
	require.NoError(b, err)
	defer func() {
		b.StopTimer()
		require.NoError(b, lim.Close(b.Context()))
	}()

	b.ReportAllocs()
	for i := range limiter.DefaultMaxKeys {
		ctx := meta.WithAttributes(b.Context(), meta.WithUserAgent(meta.String(strconv.FormatUint(i, 10))))
		_, _, err := lim.Take(ctx)
		require.NoError(b, err)
	}
	b.ResetTimer()

	for i := int(limiter.DefaultMaxKeys); b.Loop(); i++ {
		ctx := meta.WithAttributes(b.Context(), meta.WithUserAgent(meta.String(strconv.Itoa(i))))
		_, _, err := lim.Take(ctx)
		if err != nil {
			require.NoError(b, err)
		}
	}
}

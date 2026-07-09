package redis_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/cache"
	"github.com/alexfalkowski/go-service/v2/cache/config"
	drivererrors "github.com/alexfalkowski/go-service/v2/cache/driver/errors"
	cacheredis "github.com/alexfalkowski/go-service/v2/cache/driver/internal/redis"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

var _ cache.Pinger = (*cacheredis.Driver)(nil)

func TestClientClosesOnStop(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	cfg := redisConfig("redis://localhost:6379")

	d, err := cacheredis.NewDriver(lc, test.FS, cfg, nil)
	require.NoError(t, err)

	lc.RequireStart()
	lc.RequireStop()

	_, err = d.Get(t.Context(), "missing")
	require.ErrorIs(t, err, redis.ErrClosed)
}

func TestDriverPings(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	cfg := redisConfig("redis://localhost:6379")

	d, err := cacheredis.NewDriver(lc, test.FS, cfg, nil)
	require.NoError(t, err)
	t.Cleanup(lc.RequireStop)
	lc.RequireStart()

	require.NoError(t, d.Ping(t.Context()))
}

func TestParseURLDoesNotLeakCredentials(t *testing.T) {
	url := "redis://user:" + "secret%zz@localhost:6379"
	cfg := redisConfig(url)

	_, err := cacheredis.NewDriver(fxtest.NewLifecycle(t), test.FS, cfg, nil)

	require.ErrorIs(t, err, drivererrors.ErrInvalidURL)
	require.NotContains(t, err.Error(), "secret")
	require.NotContains(t, err.Error(), "redis://user")
}

func TestMetricsUnregisterOnStop(t *testing.T) {
	reader := test.EnableMetricsReader(t)
	lc := fxtest.NewLifecycle(t)
	cfg := redisConfig("redis://localhost:6379")

	_, err := cacheredis.NewDriver(lc, test.FS, cfg, nil)
	require.NoError(t, err)

	lc.RequireStart()
	require.Positive(t, redisMetricCount(t, reader))

	lc.RequireStop()
	require.Eventually(t, func() bool {
		return redisMetricCount(t, reader) == 0
	}, time.Second.Duration(), (10 * time.Millisecond).Duration())
}

func redisConfig(url string) *config.Config {
	return &config.Config{
		Options: map[string]any{
			"url": url,
		},
	}
}

func redisMetricCount(t *testing.T, reader metrics.Reader) int {
	t.Helper()

	got := &metrics.ResourceMetrics{}
	require.NoError(t, reader.Collect(t.Context(), got))

	count := 0
	for _, scope := range got.ScopeMetrics {
		for _, metric := range scope.Metrics {
			if strings.HasPrefix(metric.Name, "db.client.connections.") {
				count++
			}
		}
	}

	return count
}

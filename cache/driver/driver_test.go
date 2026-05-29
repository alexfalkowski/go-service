package driver_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/time"
	redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestNewDriver(t *testing.T) {
	tests := []struct {
		config *config.Config
		err    error
		name   string
		nil    bool
	}{
		{name: "disabled", nil: true},
		{name: "sync", config: &config.Config{Kind: "sync"}},
		{name: "unknown", config: &config.Config{Kind: "unknown"}, err: driver.ErrNotFound, nil: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := driver.NewDriver(driver.DriverParams{Config: tt.config})
			require.ErrorIs(t, err, tt.err)
			if tt.nil {
				require.Nil(t, d)
				return
			}

			require.NotNil(t, d)
		})
	}
}

func TestIsMissingError(t *testing.T) {
	cfg := &config.Config{Kind: "sync"}

	d, err := driver.NewDriver(driver.DriverParams{Config: cfg})
	require.NoError(t, err)

	_, err = d.Fetch(t.Context(), "missing")
	require.True(t, driver.IsMissingError(err))
	require.True(t, driver.IsMissingError(redis.Nil))
	require.True(t, driver.IsMissingError(errors.Prefix("wrapped", redis.Nil)))
	require.False(t, driver.IsMissingError(driver.ErrExpired))
	require.False(t, driver.IsMissingError(nil))
}

func TestRedisClientClosesOnStop(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	cfg := &config.Config{
		Kind: "redis",
		Options: map[string]any{
			"url": "redis://localhost:6379",
		},
	}

	d, err := driver.NewDriver(driver.DriverParams{
		Lifecycle: lc,
		FS:        test.FS,
		Config:    cfg,
	})
	require.NoError(t, err)

	lc.RequireStart()
	lc.RequireStop()

	_, err = d.Fetch(t.Context(), "missing")
	require.ErrorIs(t, err, redis.ErrClosed)
}

func TestRedisParseURLDoesNotLeakCredentials(t *testing.T) {
	url := "redis://user:" + "secret%zz@localhost:6379"
	cfg := &config.Config{
		Kind: "redis",
		Options: map[string]any{
			"url": url,
		},
	}

	_, err := driver.NewDriver(driver.DriverParams{
		Lifecycle: fxtest.NewLifecycle(t),
		FS:        test.FS,
		Config:    cfg,
	})

	require.ErrorIs(t, err, driver.ErrInvalidURL)
	require.NotContains(t, err.Error(), "secret")
	require.NotContains(t, err.Error(), "redis://user")
}

func TestRedisMetricsUnregisterOnStop(t *testing.T) {
	reader := test.EnableMetricsReader(t)
	lc := fxtest.NewLifecycle(t)
	cfg := &config.Config{
		Kind: "redis",
		Options: map[string]any{
			"url": "redis://localhost:6379",
		},
	}

	_, err := driver.NewDriver(driver.DriverParams{
		Lifecycle: lc,
		FS:        test.FS,
		Config:    cfg,
	})
	require.NoError(t, err)

	lc.RequireStart()
	require.Positive(t, redisMetricCount(t, reader))

	lc.RequireStop()
	require.Eventually(t, func() bool {
		return redisMetricCount(t, reader) == 0
	}, time.Second.Duration(), (10 * time.Millisecond).Duration())
}

func TestSyncDriverHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	d, err := driver.NewDriver(driver.DriverParams{Config: &config.Config{Kind: "sync"}})
	require.NoError(t, err)
	require.ErrorIs(t, d.Save(ctx, "key", "value", 0), context.Canceled)
	_, err = d.Fetch(ctx, "key")
	require.ErrorIs(t, err, context.Canceled)
	require.ErrorIs(t, d.Delete(ctx, "key"), context.Canceled)
	require.ErrorIs(t, d.Flush(ctx), context.Canceled)
}

func TestSyncDriverExpiresEntries(t *testing.T) {
	d, err := driver.NewDriver(driver.DriverParams{Config: &config.Config{Kind: "sync"}})
	require.NoError(t, err)

	require.NoError(t, d.Save(t.Context(), "key", "value", time.Nanosecond))

	var expired error
	require.Eventually(t, func() bool {
		_, expired = d.Fetch(t.Context(), "key")
		return driver.IsExpiredError(expired)
	}, time.Second.Duration(), (10 * time.Millisecond).Duration())
	require.ErrorIs(t, expired, driver.ErrExpired)

	_, err = d.Fetch(t.Context(), "key")
	require.ErrorIs(t, err, driver.ErrMissing)
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

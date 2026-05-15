package driver_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestIsMissingError(t *testing.T) {
	cfg := &config.Config{Kind: "sync"}

	d, err := driver.NewDriver(driver.DriverParams{Config: cfg})
	require.NoError(t, err)

	_, err = d.Fetch("missing")
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

	_, err = d.Fetch("missing")
	require.ErrorIs(t, err, redis.ErrClosed)
}

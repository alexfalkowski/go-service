package driver_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	"github.com/alexfalkowski/go-service/v2/errors"
	redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestIsMissingError(t *testing.T) {
	cfg := &config.Config{Kind: "sync"}

	d, err := driver.NewDriver(nil, cfg)
	require.NoError(t, err)

	_, err = d.Fetch("missing")
	require.True(t, driver.IsMissingError(err))
	require.True(t, driver.IsMissingError(redis.Nil))
	require.True(t, driver.IsMissingError(errors.Prefix("wrapped", redis.Nil)))
	require.False(t, driver.IsMissingError(driver.ErrExpired))
	require.False(t, driver.IsMissingError(nil))
}

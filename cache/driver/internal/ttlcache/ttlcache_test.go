package ttlcache_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/cache"
	"github.com/alexfalkowski/go-service/v2/cache/config"
	drivererrors "github.com/alexfalkowski/go-service/v2/cache/driver/errors"
	cachettl "github.com/alexfalkowski/go-service/v2/cache/driver/internal/ttlcache"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

var _ cache.Pinger = (*cachettl.Driver)(nil)

func TestDriverHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	d := cachettl.NewDriver(config.DefaultMaxEntries)
	require.ErrorIs(t, d.Save(ctx, "key", "value", 0), context.Canceled)
	_, err := d.Fetch(ctx, "key")
	require.ErrorIs(t, err, context.Canceled)
	require.ErrorIs(t, d.Delete(ctx, "key"), context.Canceled)
	require.ErrorIs(t, d.Flush(ctx), context.Canceled)
	require.ErrorIs(t, d.Ping(ctx), context.Canceled)
}

func TestDriverPings(t *testing.T) {
	d := cachettl.NewDriver(config.DefaultMaxEntries)

	require.NoError(t, d.Ping(t.Context()))
}

func TestDriverExpiresEntries(t *testing.T) {
	d := cachettl.NewDriver(config.DefaultMaxEntries)

	require.NoError(t, d.Save(t.Context(), "key", "value", time.Nanosecond))

	var missing error
	require.Eventually(t, func() bool {
		_, missing = d.Fetch(t.Context(), "key")
		return drivererrors.IsMissingError(missing)
	}, time.Second.Duration(), (10 * time.Millisecond).Duration())
	require.ErrorIs(t, missing, drivererrors.ErrMissing)
}

func TestDriverEvictsEntryAtCapacity(t *testing.T) {
	d := cachettl.NewDriver(2)

	require.NoError(t, d.Save(t.Context(), "first", "1", 0))
	require.NoError(t, d.Save(t.Context(), "second", "2", 0))

	require.NoError(t, d.Save(t.Context(), "third", "3", 0))

	value, err := d.Fetch(t.Context(), "third")
	require.NoError(t, err)
	require.Equal(t, "3", value)

	var misses int
	for _, key := range []string{"first", "second"} {
		_, err = d.Fetch(t.Context(), key)
		if drivererrors.IsMissingError(err) {
			misses++
		}
	}
	require.Equal(t, 1, misses)
}

func TestDriverCleansExpiredEntriesOnSave(t *testing.T) {
	d := cachettl.NewDriver(1)

	require.NoError(t, d.Save(t.Context(), "expired", "old", time.Nanosecond))
	require.Eventually(t, func() bool {
		require.NoError(t, d.Save(t.Context(), "new", "value", 0))

		_, err := d.Fetch(t.Context(), "expired")
		return errors.Is(err, drivererrors.ErrMissing)
	}, time.Second.Duration(), (10 * time.Millisecond).Duration())

	value, err := d.Fetch(t.Context(), "new")
	require.NoError(t, err)
	require.Equal(t, "value", value)
}

// Package ttlcache provides the internal ttlcache-backed cache driver.
package ttlcache

import (
	"github.com/alexfalkowski/go-service/v2/cache/driver/errors"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/jellydator/ttlcache/v3"
)

// NewDriver constructs a ttlcache-backed cache driver.
func NewDriver(maxEntries int) *Driver {
	cache := ttlcache.New(
		ttlcache.WithCapacity[string, string](uint64(max(maxEntries, 1))),
		ttlcache.WithDisableTouchOnHit[string, string](),
	)

	return &Driver{cache: cache}
}

// Driver is a ttlcache-backed cache driver.
type Driver struct {
	cache *ttlcache.Cache[string, string]
}

// Delete removes the cached key.
func (d *Driver) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	d.cache.Delete(key)

	return nil
}

// Fetch retrieves the cached value for key.
func (d *Driver) Fetch(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	item := d.cache.Get(key)
	if item == nil {
		return "", errors.ErrMissing
	}

	return item.Value(), nil
}

// Flush removes all cached keys managed by the driver.
func (d *Driver) Flush(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	d.cache.DeleteAll()

	return nil
}

// Ping verifies ttlcache is usable.
func (d *Driver) Ping(ctx context.Context) error {
	return ctx.Err()
}

// Save stores value under key for the provided lifetime.
func (d *Driver) Save(ctx context.Context, key, value string, lifetime time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	d.cache.DeleteExpired()
	d.cache.Set(key, value, ttl(lifetime).Duration())

	return nil
}

func ttl(lifetime time.Duration) time.Duration {
	if lifetime <= 0 {
		return time.Duration(ttlcache.NoTTL)
	}

	return lifetime
}

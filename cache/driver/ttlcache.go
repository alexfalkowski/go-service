package driver

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/jellydator/ttlcache/v3"
)

func newTTLCacheDriver(maxEntries int) *ttlcacheDriver {
	return &ttlcacheDriver{
		cache: ttlcache.New(
			ttlcache.WithCapacity[string, string](uint64(max(maxEntries, 1))),
			ttlcache.WithDisableTouchOnHit[string, string](),
		),
	}
}

type ttlcacheDriver struct {
	cache *ttlcache.Cache[string, string]
}

func (d *ttlcacheDriver) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	d.cache.Delete(key)

	return nil
}

func (d *ttlcacheDriver) Fetch(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	item := d.cache.Get(key)
	if item == nil {
		return "", ErrMissing
	}

	return item.Value(), nil
}

func (d *ttlcacheDriver) Flush(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	d.cache.DeleteAll()

	return nil
}

func (d *ttlcacheDriver) Save(ctx context.Context, key, value string, lifetime time.Duration) error {
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

package test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/cache"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

// Cache is a cache.Cache test double that reports hits and returns a fixed value from Fetch.
type Cache struct {
	Value string
}

// Delete removes the given key from the cache.
func (c *Cache) Delete(context.Context, string) error {
	return nil
}

// Fetch returns the cached value.
func (c *Cache) Fetch(context.Context, string) (string, error) {
	return c.Value, nil
}

// Flush clears the cache.
func (c *Cache) Flush(context.Context) error {
	return nil
}

// Save stores the value in the cache for the given TTL.
func (c *Cache) Save(context.Context, string, string, time.Duration) error {
	return nil
}

// ErrCache is a cache.Cache test double that fails fetch, delete, and save operations with ErrFailed.
//
// Flush succeeds so tests can use ErrCache in started worlds without turning cleanup into the failure under test.
type ErrCache struct{}

// Delete removes the given key from the cache.
func (*ErrCache) Delete(context.Context, string) error {
	return ErrFailed
}

// Fetch returns ErrFailed.
func (*ErrCache) Fetch(context.Context, string) (string, error) {
	return strings.Empty, ErrFailed
}

// Flush clears the cache.
func (*ErrCache) Flush(context.Context) error {
	return nil
}

// Save stores the value in the cache for the given TTL.
func (*ErrCache) Save(context.Context, string, string, time.Duration) error {
	return ErrFailed
}

func redisCache(lc di.Lifecycle) (*cache.Cache, error) {
	cfg := NewCacheConfig("redis", "snappy", "json", "redis")

	driver, err := driver.NewDriver(driver.DriverParams{
		Lifecycle: lc,
		FS:        FS,
		Config:    cfg,
	})
	if err != nil {
		return nil, err
	}

	params := cache.CacheParams{
		Config:     cfg,
		Compressor: Compressor,
		Encoder:    Encoder,
		Pool:       Pool,
		Driver:     driver,
	}

	return cache.NewCache(params), nil
}

func newWorldCache(tb testing.TB, lc di.Lifecycle, opts *worldOpts) *cache.Cache {
	tb.Helper()

	var kind *cache.Cache
	if opts.cache == nil {
		var err error
		kind, err = redisCache(lc)
		require.NoError(tb, err)
	} else {
		kind = createWorldCache(tb, lc, opts.cache)
	}

	if opts.registerCache {
		cache.Register(kind)
		tb.Cleanup(func() {
			cache.Register(nil)
		})
	}

	return kind
}

func createWorldCache(tb testing.TB, lc di.Lifecycle, opts *worldCacheOpts) *cache.Cache {
	tb.Helper()

	drv := opts.driver
	if drv == nil {
		var err error
		drv, err = driver.NewDriver(driver.DriverParams{
			Lifecycle: lc,
			FS:        FS,
			Config:    opts.config,
		})
		require.NoError(tb, err)
	}

	params := cache.CacheParams{
		Config:     opts.config,
		Compressor: Compressor,
		Encoder:    Encoder,
		Pool:       Pool,
		Driver:     drv,
	}

	return cache.NewCache(params)
}

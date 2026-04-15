package test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/v2/cache"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

// Cache is a cache.Cache test double that reports hits and returns a fixed value from Fetch.
type Cache struct {
	Value string
}

// Contains reports whether the cache contains the given key.
func (c *Cache) Contains(_ string) bool {
	return true
}

// Delete removes the given key from the cache.
func (c *Cache) Delete(_ string) error {
	return nil
}

// Fetch returns the cached value.
func (c *Cache) Fetch(_ string) (string, error) {
	return c.Value, nil
}

// FetchMulti returns cached values for the given keys.
func (c *Cache) FetchMulti(_ []string) map[string]string {
	return map[string]string{}
}

// Flush clears the cache.
func (c *Cache) Flush() error {
	return nil
}

// Save stores the value in the cache for the given TTL.
func (c *Cache) Save(_, _ string, _ time.Duration) error {
	return nil
}

// ErrCache is a cache.Cache test double that fails lookup and mutation operations with ErrFailed.
type ErrCache struct{}

// Contains reports whether the cache contains the given key.
func (*ErrCache) Contains(_ string) bool {
	return true
}

// Delete removes the given key from the cache.
func (*ErrCache) Delete(_ string) error {
	return ErrFailed
}

// Fetch returns ErrFailed.
func (*ErrCache) Fetch(_ string) (string, error) {
	return strings.Empty, ErrFailed
}

// FetchMulti returns cached values for the given keys.
func (*ErrCache) FetchMulti(_ []string) map[string]string {
	return map[string]string{}
}

// Flush clears the cache.
func (*ErrCache) Flush() error {
	return nil
}

// Save stores the value in the cache for the given TTL.
func (*ErrCache) Save(_, _ string, _ time.Duration) error {
	return ErrFailed
}

func redisCache(lc di.Lifecycle) (*cache.Cache, error) {
	cfg := NewCacheConfig("redis", "snappy", "json", "redis")

	driver, err := driver.NewDriver(FS, cfg)
	if err != nil {
		return nil, err
	}

	params := cache.CacheParams{
		Lifecycle:  lc,
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
		drv, err = driver.NewDriver(FS, opts.config)
		require.NoError(tb, err)
	}

	params := cache.CacheParams{
		Lifecycle:  lc,
		Config:     opts.config,
		Compressor: Compressor,
		Encoder:    Encoder,
		Pool:       Pool,
		Driver:     drv,
	}

	return cache.NewCache(params)
}

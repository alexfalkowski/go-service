package test

import (
	"github.com/alexfalkowski/go-service/v2/cache"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
)

// Cache is a cache implementation for tests.
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

// ErrCache is a cache implementation for tests that returns ErrFailed for most operations.
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

func redisCache(lc di.Lifecycle) *cache.Cache {
	cfg := NewCacheConfig("redis", "snappy", "json", "redis")

	driver, err := driver.NewDriver(FS, cfg)
	runtime.Must(err)

	params := cache.CacheParams{
		Lifecycle:  lc,
		Config:     cfg,
		Compressor: Compressor,
		Encoder:    Encoder,
		Pool:       Pool,
		Driver:     driver,
	}

	return cache.NewCache(params)
}

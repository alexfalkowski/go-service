package cache

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/types/ptr"
)

var cache *Cache

// Register installs the package-level cache used by generic helpers.
func Register(c *Cache) {
	cache = c
}

// Get loads a cached value for key into a new T.
//
// It returns a zero-value T and nil when caching is disabled or on cache misses.
func Get[T any](ctx context.Context, key string) (*T, error) {
	value := ptr.Zero[T]()
	if cache == nil {
		return value, nil
	}

	return value, cache.Get(ctx, key, value)
}

// Persist stores value under key with the provided TTL.
//
// It is a no-op when caching is disabled.
func Persist[T any](ctx context.Context, key string, value *T, ttl time.Duration) error {
	if cache == nil {
		return nil
	}

	return cache.Persist(ctx, key, value, ttl)
}

package cache

import (
	"context"

	"github.com/alexfalkowski/go-service/cache/cacher"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/types/ptr"
)

var cache cacher.Cache

// Register the cache.
func Register(c cacher.Cache) {
	cache = c
}

// Get a value from key.
func Get[T any](ctx context.Context, key string) (*T, bool, error) {
	value := ptr.Zero[T]()
	ok, err := cache.Get(ctx, key, value)

	return value, ok, err
}

// Persist a value to the key with a TTL.
func Persist[T any](ctx context.Context, key string, value *T, ttl time.Duration) error {
	return cache.Persist(ctx, key, value, ttl)
}

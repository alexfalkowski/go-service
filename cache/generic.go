package cache

import (
	"context"

	"github.com/alexfalkowski/go-service/cache/internal/cache"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/types/ptr"
)

var cacheable cache.Cacheable

// Register the cache.
func Register(cache cache.Cacheable) {
	cacheable = cache
}

// Get a value from key.
func Get[T any](ctx context.Context, key string) (*T, error) {
	value := ptr.Zero[T]()
	err := cacheable.Get(ctx, key, value)

	return value, err
}

// Persist a value to the key with a TTL.
func Persist[T any](ctx context.Context, key string, value *T, ttl time.Duration) error {
	return cacheable.Persist(ctx, key, value, ttl)
}

package cache

import (
	"context"

	"github.com/alexfalkowski/go-service/cache/config"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/types/ptr"
)

var cache config.Cache

// Register the cache.
func Register(c config.Cache) {
	cache = c
}

// Get a value from key.
func Get[T any](ctx context.Context, key string) (*T, error) {
	value := ptr.Zero[T]()
	err := cache.Get(ctx, key, value)

	return value, err
}

// Persist a value to the key with a TTL.
func Persist[T any](ctx context.Context, key string, value *T, ttl time.Duration) error {
	return cache.Persist(ctx, key, value, ttl)
}

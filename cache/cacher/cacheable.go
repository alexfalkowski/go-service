package cacher

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/time"
)

// Cache allows marshaling and compressing items to the cache.
type Cache interface {
	// Close the cache.
	Close(ctx context.Context) error

	// Remove a cached key.
	Remove(ctx context.Context, key string) (bool, error)

	// Get a cached value.
	Get(ctx context.Context, key string, value any) (bool, error)

	// Persist a value with key and TTL.
	Persist(ctx context.Context, key string, value any, ttl time.Duration) error
}

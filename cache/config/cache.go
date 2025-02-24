package config

import (
	"context"

	"github.com/alexfalkowski/go-service/time"
)

// Cache allows marshaling and compressing items to the cache.
type Cache interface {
	// Close the cache.
	Close(ctx context.Context) error

	// Remove a cached key.
	Remove(ctx context.Context, key string) error

	// Get a cached value.
	Get(ctx context.Context, key string, value any) error

	// Persist a value with key and TTL.
	Persist(ctx context.Context, key string, value any, ttl time.Duration) error
}

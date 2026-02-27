package cache

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/types/ptr"
)

var cache *Cache

// Register installs the package-level cache instance used by the generic helper functions.
//
// This function is primarily intended to be called by dependency injection wiring (see `Module`).
// Once registered, package-level helpers like `Get` and `Persist` will delegate to the registered
// `*Cache` instance.
//
// If c is nil, the helpers behave as if caching is disabled (they return zero values / act as no-ops).
func Register(c *Cache) {
	cache = c
}

// Get loads a cached value for key into a newly allocated value of type T and returns it.
//
// Semantics:
//   - If caching is disabled (no cache registered), Get returns a zero-value *T and a nil error.
//   - If the cache driver reports a miss/expired entry, Get returns a zero-value *T and a nil error.
//   - If a non-miss error occurs (for example decode failure or driver error), Get returns the
//     zero-value *T along with that error.
//
// The returned pointer is always non-nil.
func Get[T any](ctx context.Context, key string) (*T, error) {
	value := ptr.Zero[T]()
	if cache == nil {
		return value, nil
	}

	return value, cache.Get(ctx, key, value)
}

// Persist stores value under key with the provided TTL.
//
// If caching is disabled (no cache registered), Persist is a no-op and returns nil.
// Otherwise it delegates to the registered `*Cache`.
func Persist[T any](ctx context.Context, key string, value *T, ttl time.Duration) error {
	if cache == nil {
		return nil
	}

	return cache.Persist(ctx, key, value, ttl)
}

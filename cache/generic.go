package cache

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/ptr"
	"github.com/alexfalkowski/go-service/v2/time"
)

var cache *Cache

// Register installs the package-level cache instance used by the generic helper functions.
//
// This function is primarily intended to be called by dependency injection wiring (see [Module]).
// Once registered, package-level helpers like [Get] and [Persist] will delegate to the registered
// *[Cache] instance.
//
// If c is nil, the helpers behave as if caching is disabled: [Get] returns nil, false, nil and
// [Persist] is a no-op.
func Register(c *Cache) {
	cache = c
}

// Get loads a cached value for key into a newly allocated value of type T and reports whether a value was found.
//
// Semantics:
//   - If caching is disabled (no cache registered), [Get] returns nil, false, and a nil error.
//   - If the cache driver reports a miss/expired entry, [Get] returns a zero-value *T, false, and a nil error.
//   - If a non-miss error occurs (for example decode failure or driver error), [Get] returns the
//     zero-value *T, false, along with that error.
func Get[T any](ctx context.Context, key string) (*T, bool, error) {
	if cache == nil {
		return nil, false, nil
	}

	value := ptr.Zero[T]()
	ok, err := cache.Get(ctx, key, value)

	return value, ok, err
}

// Persist stores value under key with the provided TTL.
//
// If caching is disabled (no cache registered), [Persist] is a no-op and returns nil.
// Otherwise it delegates to the registered *[Cache].
func Persist[T any](ctx context.Context, key string, value *T, ttl time.Duration) error {
	if cache == nil {
		return nil
	}

	return cache.Persist(ctx, key, value, ttl)
}

// GetOrPersist returns the cached value for key, or produces and stores a new value of type T via fn when the key is absent.
//
// If caching is disabled (no cache registered), GetOrPersist does not call fn and returns a nil value with a
// nil error, mirroring how [Get] reports a disabled cache. Callers must tolerate a nil value rather than
// assuming a value was produced or cached.
// Concurrent in-process calls for the same key run fn once and share the produced value.
func GetOrPersist[T any](ctx context.Context, key string, ttl time.Duration, fn func() (T, error)) (*T, error) {
	if cache == nil {
		return nil, nil
	}

	value := ptr.Zero[T]()
	load := func() error {
		v, err := fn()
		if err != nil {
			return err
		}

		*value = v

		return nil
	}

	if err := cache.GetOrPersist(ctx, key, value, ttl, load); err != nil {
		return nil, err
	}

	return value, nil
}

package metrics

import (
	"github.com/alexfalkowski/go-service/v2/cache/cacher"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/time"
)

// Meter is an alias for metrics.Meter.
type Meter = metrics.Meter

// NewCache for metrics.
func NewCache(kind string, meter *metrics.Meter, cache cacher.Cache) *Cache {
	hits := meter.MustInt64Counter("cache_hits_total", "The number of hits in the cache.")
	misses := meter.MustInt64Counter("cache_misses_total", "The number of misses in the cache.")

	return &Cache{kind: kind, hits: hits, misses: misses, cache: cache}
}

// Cache for metrics.
type Cache struct {
	cache  cacher.Cache
	hits   metrics.Int64Counter
	misses metrics.Int64Counter
	kind   string
}

// Close the cache.
func (c *Cache) Close(ctx context.Context) error {
	return c.cache.Close(ctx)
}

// Remove a cached key.
func (c *Cache) Remove(ctx context.Context, key string) (bool, error) {
	return c.cache.Remove(ctx, key)
}

// Get a cached value.
func (c *Cache) Get(ctx context.Context, key string, value any) (bool, error) {
	opts := metrics.WithAttributes(attributes.String("kind", c.kind))

	ok, err := c.cache.Get(ctx, key, value)
	if err != nil {
		return ok, err
	}
	if ok {
		c.hits.Add(ctx, 1, opts)
	} else {
		c.misses.Add(ctx, 1, opts)
	}

	return ok, err
}

// Persist a value with key and TTL.
func (c *Cache) Persist(ctx context.Context, key string, value any, ttl time.Duration) error {
	return c.cache.Persist(ctx, key, value, ttl)
}

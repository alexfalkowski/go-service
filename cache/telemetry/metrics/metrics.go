package metrics

import (
	"context"

	"github.com/alexfalkowski/go-service/cache/internal/cache"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/time"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// NewCache for metrics.
func NewCache(kind string, meter *metrics.Meter, cache cache.Cacheable) *Cache {
	hits := meter.MustInt64Counter("cache_hits_total", "The number of hits in the cache.")
	misses := meter.MustInt64Counter("cache_misses_total", "The number of misses in the cache.")

	return &Cache{kind: kind, hits: hits, misses: misses, cache: cache}
}

// Cache for metrics.
type Cache struct {
	cache  cache.Cacheable
	hits   metric.Int64Counter
	misses metric.Int64Counter
	kind   string
}

// Close the cache.
func (c *Cache) Close(ctx context.Context) error {
	return c.cache.Close(ctx)
}

// Remove a cached key.
func (c *Cache) Remove(ctx context.Context, key string) error {
	return c.cache.Remove(ctx, key)
}

// Get a cached value.
func (c *Cache) Get(ctx context.Context, key string, value any) error {
	kind := attribute.Key("kind")
	opts := metric.WithAttributes(kind.String(c.kind))

	err := c.cache.Get(ctx, key, value)
	if err != nil {
		c.misses.Add(ctx, 1, opts)

		return err
	}

	c.hits.Add(ctx, 1, opts)

	return nil
}

// Persist a value with key and TTL.
func (c *Cache) Persist(ctx context.Context, key string, value any, ttl time.Duration) error {
	return c.cache.Persist(ctx, key, value, ttl)
}

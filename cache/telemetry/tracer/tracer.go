package tracer

import (
	"github.com/alexfalkowski/go-service/v2/cache/cacher"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
)

// Tracer is an alias for tracer.Tracer.
type Tracer = tracer.Tracer

// NewCache for tracer.
func NewCache(kind string, tracer *tracer.Tracer, cache cacher.Cache) *Cache {
	return &Cache{kind: kind, tracer: tracer, cache: cache}
}

// Cache for tracer.
type Cache struct {
	tracer *tracer.Tracer
	cache  cacher.Cache
	kind   string
}

// Close the cache.
func (c *Cache) Close(ctx context.Context) error {
	return c.cache.Close(ctx)
}

// Remove a cached key.
func (c *Cache) Remove(ctx context.Context, key string) (bool, error) {
	ctx, span := c.tracer.StartClient(ctx, operationName("remove"),
		attributes.String("cache.key", key),
		attributes.String("cache.kind", c.kind))
	defer span.End()

	ok, err := c.cache.Remove(ctx, key)

	span.SetAttributes(attributes.Bool("cache.found", ok))
	tracer.Error(err, span)
	tracer.Meta(ctx, span)

	return ok, err
}

// Get a cached value.
func (c *Cache) Get(ctx context.Context, key string, value any) (bool, error) {
	ctx, span := c.tracer.StartClient(ctx, operationName("get"),
		attributes.String("cache.key", key),
		attributes.String("cache.kind", c.kind))
	defer span.End()

	ok, err := c.cache.Get(ctx, key, value)

	span.SetAttributes(attributes.Bool("cache.found", ok))
	tracer.Error(err, span)
	tracer.Meta(ctx, span)

	return ok, err
}

// Persist a value with key and TTL.
func (c *Cache) Persist(ctx context.Context, key string, value any, ttl time.Duration) error {
	ctx, span := c.tracer.StartClient(ctx, operationName("persist"),
		attributes.String("cache.key", key),
		attributes.String("cache.kind", c.kind),
		attributes.String("cache.ttl", ttl.String()))
	defer span.End()

	err := c.cache.Persist(ctx, key, value, ttl)

	tracer.Error(err, span)
	tracer.Meta(ctx, span)

	return err
}

func operationName(name string) string {
	return tracer.OperationName("cache", name)
}

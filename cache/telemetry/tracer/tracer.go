package tracer

import (
	"context"

	"github.com/alexfalkowski/go-service/cache/cacher"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/time"
	"go.opentelemetry.io/otel/attribute"
)

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
	attrs := []attribute.KeyValue{
		attribute.Key("cache.key").String(key),
		attribute.Key("cache.kind").String(c.kind),
	}

	ctx, span := c.tracer.StartClient(ctx, operationName("remove"), attrs...)
	defer span.End()

	ok, err := c.cache.Remove(ctx, key)

	span.SetAttributes(attribute.Key("cache.found").Bool(ok))
	tracer.Error(err, span)
	tracer.Meta(ctx, span)

	return ok, err
}

// Get a cached value.
func (c *Cache) Get(ctx context.Context, key string, value any) (bool, error) {
	attrs := []attribute.KeyValue{
		attribute.Key("cache.key").String(key),
		attribute.Key("cache.kind").String(c.kind),
	}

	ctx, span := c.tracer.StartClient(ctx, operationName("get"), attrs...)
	defer span.End()

	ok, err := c.cache.Get(ctx, key, value)

	span.SetAttributes(attribute.Key("cache.found").Bool(ok))
	tracer.Error(err, span)
	tracer.Meta(ctx, span)

	return ok, err
}

// Persist a value with key and TTL.
func (c *Cache) Persist(ctx context.Context, key string, value any, ttl time.Duration) error {
	attrs := []attribute.KeyValue{
		attribute.Key("cache.key").String(key),
		attribute.Key("cache.kind").String(c.kind),
		attribute.Key("cache.ttl").String(ttl.String()),
	}

	ctx, span := c.tracer.StartClient(ctx, operationName("persist"), attrs...)
	defer span.End()

	err := c.cache.Persist(ctx, key, value, ttl)

	tracer.Error(err, span)
	tracer.Meta(ctx, span)

	return err
}

func operationName(name string) string {
	return tracer.OperationName("cache", name)
}

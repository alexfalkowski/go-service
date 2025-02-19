package tracer

import (
	"context"

	cache "github.com/alexfalkowski/go-service/cache/config"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/time"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// NewCache for tracer.
func NewCache(kind string, tracer *tracer.Tracer, cache cache.Cache) *Cache {
	return &Cache{kind: kind, tracer: tracer, cache: cache}
}

// Cache for tracer.
type Cache struct {
	tracer *tracer.Tracer
	cache  cache.Cache
	kind   string
}

// Close the cache.
func (c *Cache) Close(ctx context.Context) error {
	return c.cache.Close(ctx)
}

// Remove a cached key.
func (c *Cache) Remove(ctx context.Context, key string) error {
	attrs := []attribute.KeyValue{
		attribute.Key("cache.key").String(key),
		attribute.Key("cache.kind").String(c.kind),
	}

	ctx, span := c.tracer.Start(ctx, operationName("remove"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attrs...))
	defer span.End()

	ctx = tracer.WithTraceID(ctx, span)
	err := c.cache.Remove(ctx, key)

	tracer.Error(err, span)
	tracer.Meta(ctx, span)

	return err
}

// Get a cached value.
func (c *Cache) Get(ctx context.Context, key string, value any) error {
	attrs := []attribute.KeyValue{
		attribute.Key("cache.key").String(key),
		attribute.Key("cache.kind").String(c.kind),
	}

	ctx, span := c.tracer.Start(ctx, operationName("get"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attrs...))
	defer span.End()

	ctx = tracer.WithTraceID(ctx, span)
	err := c.cache.Get(ctx, key, value)

	tracer.Error(err, span)
	tracer.Meta(ctx, span)

	return err
}

// Persist a value with key and TTL.
func (c *Cache) Persist(ctx context.Context, key string, value any, ttl time.Duration) error {
	attrs := []attribute.KeyValue{
		attribute.Key("cache.key").String(key),
		attribute.Key("cache.kind").String(c.kind),
		attribute.Key("cache.ttl").String(ttl.String()),
	}

	ctx, span := c.tracer.Start(ctx, operationName("persist"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attrs...))
	defer span.End()

	ctx = tracer.WithTraceID(ctx, span)
	err := c.cache.Persist(ctx, key, value, ttl)

	tracer.Error(err, span)
	tracer.Meta(ctx, span)

	return err
}

func operationName(name string) string {
	return tracer.OperationName("cache", name)
}

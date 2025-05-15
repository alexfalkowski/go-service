package logger

import (
	"context"
	"log/slog"

	"github.com/alexfalkowski/go-service/cache/cacher"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/meta"
)

// NewCache for tracer.
func NewCache(kind string, logger *logger.Logger, cache cacher.Cache) *Cache {
	return &Cache{kind: kind, logger: logger, cache: cache}
}

// Cache for tracer.
type Cache struct {
	logger *logger.Logger
	cache  cacher.Cache
	kind   string
}

// Close the cache.
func (c *Cache) Close(ctx context.Context) error {
	return c.cache.Close(ctx)
}

// Remove a cached key.
func (c *Cache) Remove(ctx context.Context, key string) (bool, error) {
	start := time.Now()
	attrs := []slog.Attr{
		slog.String(meta.ServiceKey, c.kind),
		slog.String(meta.PathKey, key),
	}

	ok, err := c.cache.Remove(ctx, key)

	attrs = append(attrs,
		slog.String(meta.DurationKey, time.Since(start).String()),
		slog.Bool("exists", ok),
	)
	c.logger.Log(ctx, logger.NewMessage(message("remove"), err), attrs...)

	return ok, err
}

// Get a cached value.
func (c *Cache) Get(ctx context.Context, key string, value any) (bool, error) {
	start := time.Now()
	attrs := []slog.Attr{
		slog.String(meta.ServiceKey, c.kind),
		slog.String(meta.PathKey, key),
	}

	ok, err := c.cache.Get(ctx, key, value)

	attrs = append(attrs,
		slog.String(meta.DurationKey, time.Since(start).String()),
		slog.Bool("exists", ok),
	)
	c.logger.Log(ctx, logger.NewMessage(message("get"), err), attrs...)

	return ok, err
}

// Persist a value with key and TTL.
func (c *Cache) Persist(ctx context.Context, key string, value any, ttl time.Duration) error {
	start := time.Now()
	attrs := []slog.Attr{
		slog.String(meta.ServiceKey, c.kind),
		slog.String(meta.PathKey, key),
	}

	err := c.cache.Persist(ctx, key, value, ttl)

	attrs = append(attrs, slog.String(meta.DurationKey, time.Since(start).String()))
	c.logger.Log(ctx, logger.NewMessage(message("persist"), err), attrs...)

	return err
}

func message(msg string) string {
	return "cache: " + msg
}

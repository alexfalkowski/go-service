package zap

import (
	"context"
	"time"

	cache "github.com/alexfalkowski/go-service/cache/config"
	tz "github.com/alexfalkowski/go-service/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/meta"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewCache for tracer.
func NewCache(kind string, logger *zap.Logger, cache cache.Cache) *Cache {
	return &Cache{kind: kind, logger: logger, cache: cache}
}

// Cache for tracer.
type Cache struct {
	logger *zap.Logger
	cache  cache.Cache
	kind   string
}

// Close the cache.
func (c *Cache) Close(ctx context.Context) error {
	return c.cache.Close(ctx)
}

// Remove a cached key.
func (c *Cache) Remove(ctx context.Context, key string) error {
	start := time.Now()
	fields := []zapcore.Field{
		zap.String(meta.ServiceKey, c.kind),
		zap.String(meta.PathKey, key),
	}

	err := c.cache.Remove(ctx, key)
	fields = append(fields, tz.Meta(ctx)...)
	fields = append(fields, zap.Stringer(meta.DurationKey, time.Since(start)))

	tz.LogWithLogger(message("remove"), err, c.logger, fields...)

	return err
}

// Get a cached value.
func (c *Cache) Get(ctx context.Context, key string, value any) error {
	start := time.Now()
	fields := []zapcore.Field{
		zap.String(meta.ServiceKey, c.kind),
		zap.String(meta.PathKey, key),
	}

	err := c.cache.Get(ctx, key, value)
	fields = append(fields, tz.Meta(ctx)...)
	fields = append(fields, zap.Stringer(meta.DurationKey, time.Since(start)))

	tz.LogWithLogger(message("get"), err, c.logger, fields...)

	return err
}

// Persist a value with key and TTL.
func (c *Cache) Persist(ctx context.Context, key string, value any, ttl time.Duration) error {
	start := time.Now()
	fields := []zapcore.Field{
		zap.String(meta.ServiceKey, c.kind),
		zap.String(meta.PathKey, key),
	}

	err := c.cache.Persist(ctx, key, value, ttl)
	fields = append(fields, tz.Meta(ctx)...)
	fields = append(fields, zap.Stringer(meta.DurationKey, time.Since(start)))

	tz.LogWithLogger(message("persist"), err, c.logger, fields...)

	return err
}

func message(msg string) string {
	return "cache: " + msg
}

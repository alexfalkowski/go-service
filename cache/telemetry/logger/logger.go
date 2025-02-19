package logger

import (
	"context"

	cache "github.com/alexfalkowski/go-service/cache/config"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/meta"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewCache for tracer.
func NewCache(kind string, logger *logger.Logger, cache cache.Cache) *Cache {
	return &Cache{kind: kind, logger: logger, cache: cache}
}

// Cache for tracer.
type Cache struct {
	logger *logger.Logger
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
	fields = append(fields, logger.Meta(ctx)...)
	fields = append(fields, zap.Stringer(meta.DurationKey, time.Since(start)))

	c.logger.Log(message("remove"), err, fields...)

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
	fields = append(fields, logger.Meta(ctx)...)
	fields = append(fields, zap.Stringer(meta.DurationKey, time.Since(start)))

	c.logger.Log(message("get"), err, fields...)

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
	fields = append(fields, logger.Meta(ctx)...)
	fields = append(fields, zap.Stringer(meta.DurationKey, time.Since(start)))

	c.logger.Log(message("persist"), err, fields...)

	return err
}

func message(msg string) string {
	return "cache: " + msg
}

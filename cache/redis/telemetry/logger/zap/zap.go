package zap

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/cache/redis/client"
	"github.com/alexfalkowski/go-service/meta"
	stime "github.com/alexfalkowski/go-service/time"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	redisDuration  = "redis.duration"
	redisStartTime = "redis.start_time"
	redisDeadline  = "redis.deadline"
	component      = "component"
	redisComponent = "redis"
	kind           = "redis"
)

// NewClient for zap.
func NewClient(logger *zap.Logger, client client.Client) *Client {
	return &Client{logger: logger, client: client}
}

// Client for zap.
type Client struct {
	logger *zap.Logger
	client client.Client
}

//nolint:dupl
func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) *redis.StatusCmd {
	start := time.Now().UTC()
	cmd := c.client.Set(ctx, key, value, ttl)
	fields := []zapcore.Field{
		zap.Int64(redisDuration, stime.ToMilliseconds(time.Since(start))),
		zap.String(redisStartTime, start.Format(time.RFC3339)),
		zap.String("span.kind", kind),
		zap.String(component, redisComponent),
		zap.String("redis.key", key),
		zap.Duration("redis.ttl", ttl),
	}

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(redisDeadline, d.UTC().Format(time.RFC3339)))
	}

	if err := cmd.Err(); err != nil {
		fields = append(fields, zap.Error(err))
		c.logger.Error("finished call with error", fields...)
	} else {
		c.logger.Info("finished call with success", fields...)
	}

	return cmd
}

//nolint:dupl
func (c *Client) SetXX(ctx context.Context, key string, value any, ttl time.Duration) *redis.BoolCmd {
	start := time.Now().UTC()
	cmd := c.client.SetXX(ctx, key, value, ttl)
	fields := []zapcore.Field{
		zap.Int64(redisDuration, stime.ToMilliseconds(time.Since(start))),
		zap.String(redisStartTime, start.Format(time.RFC3339)),
		zap.String("span.kind", kind),
		zap.String(component, redisComponent),
		zap.String("redis.key", key),
		zap.Duration("redis.ttl", ttl),
	}

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(redisDeadline, d.UTC().Format(time.RFC3339)))
	}

	if err := cmd.Err(); err != nil {
		fields = append(fields, zap.Error(err))
		c.logger.Error("finished call with error", fields...)
	} else {
		c.logger.Info("finished call with success", fields...)
	}

	return cmd
}

//nolint:dupl
func (c *Client) SetNX(ctx context.Context, key string, value any, ttl time.Duration) *redis.BoolCmd {
	start := time.Now().UTC()
	cmd := c.client.SetNX(ctx, key, value, ttl)
	fields := []zapcore.Field{
		zap.Int64(redisDuration, stime.ToMilliseconds(time.Since(start))),
		zap.String(redisStartTime, start.Format(time.RFC3339)),
		zap.String("span.kind", kind),
		zap.String(component, redisComponent),
		zap.String("redis.key", key),
		zap.Duration("redis.ttl", ttl),
	}

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(redisDeadline, d.UTC().Format(time.RFC3339)))
	}

	if err := cmd.Err(); err != nil {
		fields = append(fields, zap.Error(err))
		c.logger.Error("finished call with error", fields...)
	} else {
		c.logger.Info("finished call with success", fields...)
	}

	return cmd
}

//nolint:dupl
func (c *Client) Get(ctx context.Context, key string) *redis.StringCmd {
	start := time.Now().UTC()
	cmd := c.client.Get(ctx, key)
	fields := []zapcore.Field{
		zap.Int64(redisDuration, stime.ToMilliseconds(time.Since(start))),
		zap.String(redisStartTime, start.Format(time.RFC3339)),
		zap.String("span.kind", kind),
		zap.String(component, redisComponent),
		zap.String("redis.key", key),
	}

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(redisDeadline, d.UTC().Format(time.RFC3339)))
	}

	if err := cmd.Err(); err != nil {
		fields = append(fields, zap.Error(err))
		c.logger.Error("finished call with error", fields...)
	} else {
		c.logger.Info("finished call with success", fields...)
	}

	return cmd
}

//nolint:dupl
func (c *Client) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	start := time.Now().UTC()
	cmd := c.client.Del(ctx, keys...)
	fields := []zapcore.Field{
		zap.Int64(redisDuration, stime.ToMilliseconds(time.Since(start))),
		zap.String(redisStartTime, start.Format(time.RFC3339)),
		zap.String("span.kind", kind),
		zap.String(component, redisComponent),
		zap.Strings("redis.keys", keys),
	}

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(redisDeadline, d.UTC().Format(time.RFC3339)))
	}

	if err := cmd.Err(); err != nil {
		fields = append(fields, zap.Error(err))
		c.logger.Error("finished call with error", fields...)
	} else {
		c.logger.Info("finished call with success", fields...)
	}

	return cmd
}

//nolint:dupl
func (c *Client) Incr(ctx context.Context, key string) *redis.IntCmd {
	start := time.Now().UTC()
	cmd := c.client.Incr(ctx, key)
	fields := []zapcore.Field{
		zap.Int64(redisDuration, stime.ToMilliseconds(time.Since(start))),
		zap.String(redisStartTime, start.Format(time.RFC3339)),
		zap.String("span.kind", kind),
		zap.String(component, redisComponent),
		zap.String("redis.key", key),
	}

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(redisDeadline, d.UTC().Format(time.RFC3339)))
	}

	if err := cmd.Err(); err != nil {
		fields = append(fields, zap.Error(err))
		c.logger.Error("finished call with error", fields...)
	} else {
		c.logger.Info("finished call with success", fields...)
	}

	return cmd
}

func (c *Client) Ping(ctx context.Context) *redis.StatusCmd {
	return c.client.Ping(ctx)
}

func (c *Client) Close() error {
	return c.client.Close()
}

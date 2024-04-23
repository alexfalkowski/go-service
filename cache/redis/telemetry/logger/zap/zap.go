package zap

import (
	"context"
	"time"

	gr "github.com/alexfalkowski/go-service/redis"
	tz "github.com/alexfalkowski/go-service/telemetry/logger/zap"
	st "github.com/alexfalkowski/go-service/time"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	service = "redis"
)

// NewClient for zap.
func NewClient(logger *zap.Logger, client gr.Client) *Client {
	return &Client{logger: logger, client: client}
}

// Client for zap.
type Client struct {
	logger *zap.Logger
	client gr.Client
}

//nolint:dupl
func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) *redis.StatusCmd {
	start := time.Now()
	cmd := c.client.Set(ctx, key, value, ttl)
	fields := []zapcore.Field{
		zap.Int64(tm.DurationKey, st.ToMilliseconds(time.Since(start))),
		zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
		zap.String(tm.ServiceKey, service),
		zap.String(tm.PathKey, key),
	}

	fields = append(fields, tz.Meta(ctx)...)

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(tm.DeadlineKey, d.Format(time.RFC3339)))
	}

	tz.LogWithLogger(message("client set"), cmd.Err(), c.logger, fields...)

	return cmd
}

//nolint:dupl
func (c *Client) SetXX(ctx context.Context, key string, value any, ttl time.Duration) *redis.BoolCmd {
	start := time.Now()
	cmd := c.client.SetXX(ctx, key, value, ttl)
	fields := []zapcore.Field{
		zap.Int64(tm.DurationKey, st.ToMilliseconds(time.Since(start))),
		zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
		zap.String(tm.ServiceKey, service),
		zap.String(tm.PathKey, key),
	}

	fields = append(fields, tz.Meta(ctx)...)

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(tm.DeadlineKey, d.Format(time.RFC3339)))
	}

	tz.LogWithLogger(message("client setxx"), cmd.Err(), c.logger, fields...)

	return cmd
}

//nolint:dupl
func (c *Client) SetNX(ctx context.Context, key string, value any, ttl time.Duration) *redis.BoolCmd {
	start := time.Now()
	cmd := c.client.SetNX(ctx, key, value, ttl)
	fields := []zapcore.Field{
		zap.Int64(tm.DurationKey, st.ToMilliseconds(time.Since(start))),
		zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
		zap.String(tm.ServiceKey, service),
		zap.String(tm.PathKey, key),
	}

	fields = append(fields, tz.Meta(ctx)...)

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(tm.DeadlineKey, d.Format(time.RFC3339)))
	}

	tz.LogWithLogger(message("client setnx"), cmd.Err(), c.logger, fields...)

	return cmd
}

//nolint:dupl
func (c *Client) Get(ctx context.Context, key string) *redis.StringCmd {
	start := time.Now()
	cmd := c.client.Get(ctx, key)
	fields := []zapcore.Field{
		zap.Int64(tm.DurationKey, st.ToMilliseconds(time.Since(start))),
		zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
		zap.String(tm.ServiceKey, service),
		zap.String(tm.PathKey, key),
	}

	fields = append(fields, tz.Meta(ctx)...)

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(tm.DeadlineKey, d.Format(time.RFC3339)))
	}

	tz.LogWithLogger(message("client get"), cmd.Err(), c.logger, fields...)

	return cmd
}

func (c *Client) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	start := time.Now()
	cmd := c.client.Del(ctx, keys...)
	fields := []zapcore.Field{
		zap.Int64(tm.DurationKey, st.ToMilliseconds(time.Since(start))),
		zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
		zap.String(tm.ServiceKey, service),
		zap.Strings(tm.PathKey, keys),
	}

	fields = append(fields, tz.Meta(ctx)...)

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(tm.DeadlineKey, d.Format(time.RFC3339)))
	}

	tz.LogWithLogger(message("client del"), cmd.Err(), c.logger, fields...)

	return cmd
}

//nolint:dupl
func (c *Client) Incr(ctx context.Context, key string) *redis.IntCmd {
	start := time.Now()
	cmd := c.client.Incr(ctx, key)
	fields := []zapcore.Field{
		zap.Int64(tm.DurationKey, st.ToMilliseconds(time.Since(start))),
		zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
		zap.String(tm.ServiceKey, service),
		zap.String(tm.PathKey, key),
	}

	fields = append(fields, tz.Meta(ctx)...)

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(tm.DeadlineKey, d.Format(time.RFC3339)))
	}

	tz.LogWithLogger(message("client incr"), cmd.Err(), c.logger, fields...)

	return cmd
}

func (c *Client) Ping(ctx context.Context) *redis.StatusCmd {
	return c.client.Ping(ctx)
}

func (c *Client) Close() error {
	return c.client.Close()
}

func message(msg string) string {
	return "redis: " + msg
}

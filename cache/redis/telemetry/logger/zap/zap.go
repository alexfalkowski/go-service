package zap

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	gr "github.com/alexfalkowski/go-service/redis"
	stime "github.com/alexfalkowski/go-service/time"
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
		zap.Int64(tm.DurationKey, stime.ToMilliseconds(time.Since(start))),
		zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
		zap.String(tm.ServiceKey, service),
		zap.String(tm.PathKey, key),
	}

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(tm.DeadlineKey, d.Format(time.RFC3339)))
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
	start := time.Now()
	cmd := c.client.SetXX(ctx, key, value, ttl)
	fields := []zapcore.Field{
		zap.Int64(tm.DurationKey, stime.ToMilliseconds(time.Since(start))),
		zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
		zap.String(tm.ServiceKey, service),
		zap.String(tm.PathKey, key),
	}

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(tm.DeadlineKey, d.Format(time.RFC3339)))
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
	start := time.Now()
	cmd := c.client.SetNX(ctx, key, value, ttl)
	fields := []zapcore.Field{
		zap.Int64(tm.DurationKey, stime.ToMilliseconds(time.Since(start))),
		zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
		zap.String(tm.ServiceKey, service),
		zap.String(tm.PathKey, key),
	}

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(tm.DeadlineKey, d.Format(time.RFC3339)))
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
	start := time.Now()
	cmd := c.client.Get(ctx, key)
	fields := []zapcore.Field{
		zap.Int64(tm.DurationKey, stime.ToMilliseconds(time.Since(start))),
		zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
		zap.String(tm.ServiceKey, service),
		zap.String(tm.PathKey, key),
	}

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(tm.DeadlineKey, d.Format(time.RFC3339)))
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
	start := time.Now()
	cmd := c.client.Del(ctx, keys...)
	fields := []zapcore.Field{
		zap.Int64(tm.DurationKey, stime.ToMilliseconds(time.Since(start))),
		zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
		zap.String(tm.ServiceKey, service),
		zap.Strings(tm.PathKey, keys),
	}

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(tm.DeadlineKey, d.Format(time.RFC3339)))
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
	start := time.Now()
	cmd := c.client.Incr(ctx, key)
	fields := []zapcore.Field{
		zap.Int64(tm.DurationKey, stime.ToMilliseconds(time.Since(start))),
		zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
		zap.String(tm.ServiceKey, service),
		zap.String(tm.PathKey, key),
	}

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(tm.DeadlineKey, d.Format(time.RFC3339)))
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

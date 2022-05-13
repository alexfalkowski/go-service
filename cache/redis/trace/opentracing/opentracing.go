package opentracing

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/cache/redis/client"
	"github.com/alexfalkowski/go-service/meta"
	stime "github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/trace/opentracing"
	"github.com/alexfalkowski/go-service/version"
	"github.com/go-redis/redis/v8"
	otr "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/fx"
)

const (
	redisDuration        = "redis.duration"
	redisStartTime       = "redis.start_time"
	redisRequestDeadline = "redis.request.deadline"
	component            = "component"
	redisComponent       = "redis"
)

// TracerParams for opentracing.
type TracerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *opentracing.Config
	Version   version.Version
}

// NewTracer for opentracing.
func NewTracer(params TracerParams) (Tracer, error) {
	return opentracing.NewTracer(opentracing.TracerParams{Lifecycle: params.Lifecycle, Name: "redis", Config: params.Config, Version: params.Version})
}

// Tracer for opentracing.
type Tracer otr.Tracer

// StartSpanFromContext for opentracing.
func StartSpanFromContext(ctx context.Context, tracer Tracer, operation, method string, opts ...otr.StartSpanOption) (context.Context, otr.Span) {
	return opentracing.StartSpanFromContext(ctx, tracer, "redis", operation, method, opts...)
}

// NewClient for opentracing.
func NewClient(tracer Tracer, client client.Client) *Client {
	return &Client{tracer: tracer, client: client}
}

// Client for opentracing.
type Client struct {
	tracer Tracer
	client client.Client
}

func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) *redis.StatusCmd {
	start := time.Now().UTC()
	opts := []otr.StartSpanOption{
		otr.Tag{Key: redisStartTime, Value: start.Format(time.RFC3339)},
		otr.Tag{Key: component, Value: redisComponent},
		otr.Tag{Key: "redis.key", Value: key},
		otr.Tag{Key: "redis.ttl", Value: ttl},
	}

	ctx, span := StartSpanFromContext(ctx, c.tracer, "client", "set", opts...)
	defer span.Finish()

	if d, ok := ctx.Deadline(); ok {
		span.SetTag(redisRequestDeadline, d.UTC().Format(time.RFC3339))
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	cmd := c.client.Set(ctx, key, value, ttl)
	if err := cmd.Err(); err != nil {
		setError(span, err)
	}

	span.SetTag(redisDuration, stime.ToMilliseconds(time.Since(start)))

	return cmd
}

func (c *Client) SetXX(ctx context.Context, key string, value any, ttl time.Duration) *redis.BoolCmd {
	start := time.Now().UTC()
	opts := []otr.StartSpanOption{
		otr.Tag{Key: redisStartTime, Value: start.Format(time.RFC3339)},
		otr.Tag{Key: component, Value: redisComponent},
		otr.Tag{Key: "redis.key", Value: key},
		otr.Tag{Key: "redis.ttl", Value: ttl},
	}

	ctx, span := StartSpanFromContext(ctx, c.tracer, "client", "set", opts...)
	defer span.Finish()

	if d, ok := ctx.Deadline(); ok {
		span.SetTag(redisRequestDeadline, d.UTC().Format(time.RFC3339))
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	cmd := c.client.SetXX(ctx, key, value, ttl)
	if err := cmd.Err(); err != nil {
		setError(span, err)
	}

	span.SetTag(redisDuration, stime.ToMilliseconds(time.Since(start)))

	return cmd
}

func (c *Client) SetNX(ctx context.Context, key string, value any, ttl time.Duration) *redis.BoolCmd {
	start := time.Now().UTC()
	opts := []otr.StartSpanOption{
		otr.Tag{Key: redisStartTime, Value: start.Format(time.RFC3339)},
		otr.Tag{Key: component, Value: redisComponent},
		otr.Tag{Key: "redis.key", Value: key},
		otr.Tag{Key: "redis.ttl", Value: ttl},
	}

	ctx, span := StartSpanFromContext(ctx, c.tracer, "client", "set", opts...)
	defer span.Finish()

	if d, ok := ctx.Deadline(); ok {
		span.SetTag(redisRequestDeadline, d.UTC().Format(time.RFC3339))
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	cmd := c.client.SetNX(ctx, key, value, ttl)
	if err := cmd.Err(); err != nil {
		setError(span, err)
	}

	span.SetTag(redisDuration, stime.ToMilliseconds(time.Since(start)))

	return cmd
}

func (c *Client) Get(ctx context.Context, key string) *redis.StringCmd {
	start := time.Now().UTC()
	opts := []otr.StartSpanOption{
		otr.Tag{Key: redisStartTime, Value: start.Format(time.RFC3339)},
		otr.Tag{Key: component, Value: redisComponent},
		otr.Tag{Key: "redis.key", Value: key},
	}

	ctx, span := StartSpanFromContext(ctx, c.tracer, "client", "get", opts...)
	defer span.Finish()

	if d, ok := ctx.Deadline(); ok {
		span.SetTag(redisRequestDeadline, d.UTC().Format(time.RFC3339))
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	cmd := c.client.Get(ctx, key)
	if err := cmd.Err(); err != nil {
		setError(span, err)
	}

	span.SetTag(redisDuration, stime.ToMilliseconds(time.Since(start)))

	return cmd
}

func (c *Client) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	start := time.Now().UTC()
	opts := []otr.StartSpanOption{
		otr.Tag{Key: redisStartTime, Value: start.Format(time.RFC3339)},
		otr.Tag{Key: component, Value: redisComponent},
		otr.Tag{Key: "redis.keys", Value: keys},
	}

	ctx, span := StartSpanFromContext(ctx, c.tracer, "client", "del", opts...)
	defer span.Finish()

	if d, ok := ctx.Deadline(); ok {
		span.SetTag(redisRequestDeadline, d.UTC().Format(time.RFC3339))
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	cmd := c.client.Del(ctx, keys...)
	if err := cmd.Err(); err != nil {
		setError(span, err)
	}

	span.SetTag(redisDuration, stime.ToMilliseconds(time.Since(start)))

	return cmd
}

func (c *Client) Ping(ctx context.Context) *redis.StatusCmd {
	return c.client.Ping(ctx)
}

func (c *Client) Close() error {
	return c.client.Close()
}

func setError(span otr.Span, err error) {
	ext.Error.Set(span, true)
	span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
}

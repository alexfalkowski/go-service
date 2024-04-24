package tracer

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	gr "github.com/alexfalkowski/go-service/redis"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
)

// NewClient for tracer.
func NewClient(tracer trace.Tracer, client gr.Client) *Client {
	return &Client{tracer: tracer, client: client}
}

// Client for tracer.
type Client struct {
	tracer trace.Tracer
	client gr.Client
}

func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) *redis.StatusCmd {
	attrs := []attribute.KeyValue{
		semconv.DBSystemRedis,
		attribute.Key("db.redis.key").String(key),
		attribute.Key("db.redis.ttl").Int64(int64(ttl)),
	}

	ctx, span := c.tracer.Start(ctx, operationName("client set"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attrs...))
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))
	cmd := c.client.Set(ctx, key, value, ttl)

	tracer.Error(cmd.Err(), span)
	tracer.Meta(ctx, span)

	return cmd
}

func (c *Client) SetXX(ctx context.Context, key string, value any, ttl time.Duration) *redis.BoolCmd {
	attrs := []attribute.KeyValue{
		semconv.DBSystemRedis,
		attribute.Key("db.redis.key").String(key),
		attribute.Key("db.redis.ttl").Int64(int64(ttl)),
	}

	ctx, span := c.tracer.Start(ctx, operationName("client setxx"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attrs...))
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))
	cmd := c.client.SetXX(ctx, key, value, ttl)

	tracer.Error(cmd.Err(), span)
	tracer.Meta(ctx, span)

	return cmd
}

func (c *Client) SetNX(ctx context.Context, key string, value any, ttl time.Duration) *redis.BoolCmd {
	attrs := []attribute.KeyValue{
		semconv.DBSystemRedis,
		attribute.Key("db.redis.key").String(key),
		attribute.Key("db.redis.ttl").Int64(int64(ttl)),
	}

	ctx, span := c.tracer.Start(ctx, operationName("client setnx"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attrs...))
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))
	cmd := c.client.SetNX(ctx, key, value, ttl)

	tracer.Error(cmd.Err(), span)
	tracer.Meta(ctx, span)

	return cmd
}

func (c *Client) Get(ctx context.Context, key string) *redis.StringCmd {
	attrs := []attribute.KeyValue{
		semconv.DBSystemRedis,
		attribute.Key("db.redis.key").String(key),
	}

	ctx, span := c.tracer.Start(ctx, operationName("client get"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attrs...))
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))
	cmd := c.client.Get(ctx, key)

	tracer.Error(cmd.Err(), span)
	tracer.Meta(ctx, span)

	return cmd
}

func (c *Client) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	attrs := []attribute.KeyValue{
		semconv.DBSystemRedis,
		attribute.Key("db.redis.keys").StringSlice(keys),
	}

	ctx, span := c.tracer.Start(ctx, operationName("client del"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attrs...))
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))
	cmd := c.client.Del(ctx, keys...)

	tracer.Error(cmd.Err(), span)
	tracer.Meta(ctx, span)

	return cmd
}

func (c *Client) Incr(ctx context.Context, key string) *redis.IntCmd {
	attrs := []attribute.KeyValue{
		semconv.DBSystemRedis,
		attribute.Key("db.redis.key").String(key),
	}

	ctx, span := c.tracer.Start(ctx, operationName("client incr"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attrs...))
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))
	cmd := c.client.Incr(ctx, key)

	tracer.Error(cmd.Err(), span)
	tracer.Meta(ctx, span)

	return cmd
}

func (c *Client) Ping(ctx context.Context) *redis.StatusCmd {
	return c.client.Ping(ctx)
}

func (c *Client) Close() error {
	return c.client.Close()
}

func operationName(name string) string {
	return tracer.OperationName("redis", name)
}

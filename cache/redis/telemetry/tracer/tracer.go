package tracer

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	gr "github.com/alexfalkowski/go-service/redis"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
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

//nolint:dupl
func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) *redis.StatusCmd {
	operationName := "client set"
	attrs := []attribute.KeyValue{
		semconv.DBSystemRedis,
		attribute.Key("db.redis.key").String(key),
		attribute.Key("db.redis.ttl").Int64(int64(ttl)),
	}

	ctx, span := c.tracer.Start(
		ctx,
		operationName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))

	cmd := c.client.Set(ctx, key, value, ttl)
	if err := cmd.Err(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	for k, v := range meta.Strings(ctx) {
		span.SetAttributes(attribute.Key(k).String(v))
	}

	return cmd
}

//nolint:dupl
func (c *Client) SetXX(ctx context.Context, key string, value any, ttl time.Duration) *redis.BoolCmd {
	operationName := "client setxx"
	attrs := []attribute.KeyValue{
		semconv.DBSystemRedis,
		attribute.Key("db.redis.key").String(key),
		attribute.Key("db.redis.ttl").Int64(int64(ttl)),
	}

	ctx, span := c.tracer.Start(
		ctx,
		operationName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))

	cmd := c.client.SetXX(ctx, key, value, ttl)
	if err := cmd.Err(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	for k, v := range meta.Strings(ctx) {
		span.SetAttributes(attribute.Key(k).String(v))
	}

	return cmd
}

//nolint:dupl
func (c *Client) SetNX(ctx context.Context, key string, value any, ttl time.Duration) *redis.BoolCmd {
	operationName := "client setnx"
	attrs := []attribute.KeyValue{
		semconv.DBSystemRedis,
		attribute.Key("db.redis.key").String(key),
		attribute.Key("db.redis.ttl").Int64(int64(ttl)),
	}

	ctx, span := c.tracer.Start(
		ctx,
		operationName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))

	cmd := c.client.SetNX(ctx, key, value, ttl)
	if err := cmd.Err(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	for k, v := range meta.Strings(ctx) {
		span.SetAttributes(attribute.Key(k).String(v))
	}

	return cmd
}

//nolint:dupl
func (c *Client) Get(ctx context.Context, key string) *redis.StringCmd {
	operationName := "client get"
	attrs := []attribute.KeyValue{
		semconv.DBSystemRedis,
		attribute.Key("db.redis.key").String(key),
	}

	ctx, span := c.tracer.Start(
		ctx,
		operationName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))

	cmd := c.client.Get(ctx, key)
	if err := cmd.Err(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	for k, v := range meta.Strings(ctx) {
		span.SetAttributes(attribute.Key(k).String(v))
	}

	return cmd
}

func (c *Client) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	operationName := "client del"
	attrs := []attribute.KeyValue{
		semconv.DBSystemRedis,
		attribute.Key("db.redis.keys").StringSlice(keys),
	}

	ctx, span := c.tracer.Start(
		ctx,
		operationName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))

	cmd := c.client.Del(ctx, keys...)
	if err := cmd.Err(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	for k, v := range meta.Strings(ctx) {
		span.SetAttributes(attribute.Key(k).String(v))
	}

	return cmd
}

//nolint:dupl
func (c *Client) Incr(ctx context.Context, key string) *redis.IntCmd {
	operationName := "client incr"
	attrs := []attribute.KeyValue{
		semconv.DBSystemRedis,
		attribute.Key("db.redis.key").String(key),
	}

	ctx, span := c.tracer.Start(
		ctx,
		operationName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))

	cmd := c.client.Incr(ctx, key)
	if err := cmd.Err(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	for k, v := range meta.Strings(ctx) {
		span.SetAttributes(attribute.Key(k).String(v))
	}

	return cmd
}

func (c *Client) Ping(ctx context.Context) *redis.StatusCmd {
	return c.client.Ping(ctx)
}

func (c *Client) Close() error {
	return c.client.Close()
}

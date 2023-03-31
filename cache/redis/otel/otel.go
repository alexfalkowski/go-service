package otel

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/cache/redis/client"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/otel"
	"github.com/alexfalkowski/go-service/version"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// TracerParams for otel.
type TracerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *otel.Config
	Version   version.Version
}

// NewTracer for otel.
func NewTracer(params TracerParams) (Tracer, error) {
	return otel.NewTracer(otel.TracerParams{Lifecycle: params.Lifecycle, Name: "redis", Config: params.Config, Version: params.Version})
}

// Tracer for otel.
type Tracer trace.Tracer

// NewClient for otel.
func NewClient(tracer Tracer, client client.Client) *Client {
	return &Client{tracer: tracer, client: client}
}

// Client for otel.
type Client struct {
	tracer Tracer
	client client.Client
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

	cmd := c.client.Set(ctx, key, value, ttl)
	if err := cmd.Err(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	for k, v := range meta.Attributes(ctx) {
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

	cmd := c.client.SetXX(ctx, key, value, ttl)
	if err := cmd.Err(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	for k, v := range meta.Attributes(ctx) {
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

	cmd := c.client.SetNX(ctx, key, value, ttl)
	if err := cmd.Err(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetAttributes(attribute.Key(k).String(v))
	}

	return cmd
}

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

	cmd := c.client.Get(ctx, key)
	if err := cmd.Err(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	for k, v := range meta.Attributes(ctx) {
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

	cmd := c.client.Del(ctx, keys...)
	if err := cmd.Err(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetAttributes(attribute.Key(k).String(v))
	}

	return cmd
}

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

	cmd := c.client.Incr(ctx, key)
	if err := cmd.Err(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	for k, v := range meta.Attributes(ctx) {
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

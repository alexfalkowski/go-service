package telemetry

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/cache/redis/client"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/version"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// TracerParams for telemetry.
type TracerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *telemetry.Config
	Version   version.Version
}

// NewTracer for telemetry.
func NewTracer(params TracerParams) (Tracer, error) {
	return telemetry.NewTracer(telemetry.TracerParams{Lifecycle: params.Lifecycle, Name: "redis", Config: params.Config, Version: params.Version})
}

// Tracer for telemetry.
type Tracer trace.Tracer

// NewTracerClient for telemetry.
func NewTracerClient(tracer Tracer, client client.Client) *TracerClient {
	return &TracerClient{tracer: tracer, client: client}
}

// TracerClient for telemetry.
type TracerClient struct {
	tracer Tracer
	client client.Client
}

//nolint:dupl
func (c *TracerClient) Set(ctx context.Context, key string, value any, ttl time.Duration) *redis.StatusCmd {
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
func (c *TracerClient) SetXX(ctx context.Context, key string, value any, ttl time.Duration) *redis.BoolCmd {
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
func (c *TracerClient) SetNX(ctx context.Context, key string, value any, ttl time.Duration) *redis.BoolCmd {
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

func (c *TracerClient) Get(ctx context.Context, key string) *redis.StringCmd {
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

func (c *TracerClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
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

func (c *TracerClient) Incr(ctx context.Context, key string) *redis.IntCmd {
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

func (c *TracerClient) Ping(ctx context.Context) *redis.StatusCmd {
	return c.client.Ping(ctx)
}

func (c *TracerClient) Close() error {
	return c.client.Close()
}

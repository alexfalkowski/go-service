package tracer

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/strings"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	otlp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// OperationName for tracer.
func OperationName(area, name string) string {
	return strings.Join(": ", area, name)
}

// Register for tracer.
func Register() {
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

// Params for tracer.
type Params struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Config      *Config
	ID          env.ID
	Name        env.Name
	Version     env.Version
	Environment env.Environment
}

// NewTracer for tracer.
func NewTracer(params Params) *Tracer {
	if !IsEnabled(params.Config) {
		return nil
	}

	client := otlp.NewClient(otlp.WithEndpointURL(params.Config.URL), otlp.WithHeaders(params.Config.Headers))
	exporter := otlptrace.NewUnstarted(client)
	attrs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.HostID(params.ID.String()),
		semconv.ServiceName(params.Name.String()),
		semconv.ServiceVersion(params.Version.String()),
		semconv.DeploymentEnvironmentName(params.Environment.String()),
	)
	provider := sdktrace.NewTracerProvider(sdktrace.WithResource(attrs), sdktrace.WithBatcher(exporter))

	otel.SetTracerProvider(provider)

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return prefix(exporter.Start(ctx))
		},
		OnStop: func(ctx context.Context) error {
			_ = provider.Shutdown(ctx)
			_ = exporter.Shutdown(ctx)

			return nil
		},
	})

	return &Tracer{provider.Tracer(params.Name.String())}
}

// Tracer using otel.
type Tracer struct {
	trace.Tracer
}

// StartClient starts a new client span.
//
//nolint:spancheck
func (t *Tracer) StartClient(ctx context.Context, spanName string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	ctx, span := t.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attrs...))
	ctx = WithTraceID(ctx, span)

	return ctx, span
}

// StartServer starts a new server span.
//
//nolint:spancheck
func (t *Tracer) StartServer(ctx context.Context, spanName string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	ctx = trace.ContextWithRemoteSpanContext(ctx, trace.SpanContextFromContext(ctx))
	ctx, span := t.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(attrs...))
	ctx = WithTraceID(ctx, span)

	return ctx, span
}

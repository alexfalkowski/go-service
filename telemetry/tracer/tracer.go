package tracer

import (
	"context"

	"github.com/alexfalkowski/go-service/env"
	se "github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	otlp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// OperationName for tracer.
func OperationName(area, name string) string {
	return area + ": " + name
}

// Register for tracer.
func Register() {
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

// Params for tracer.
type Params struct {
	fx.In
	Lifecycle   fx.Lifecycle
	FileSystem  os.FileSystem
	Config      *Config
	Logger      *logger.Logger
	Environment env.Environment
	Version     env.Version
	Name        env.Name
}

// NewTracer for tracer.
func NewTracer(params Params) (*Tracer, error) {
	if !IsEnabled(params.Config) {
		return &Tracer{noop.Tracer{}}, nil
	}

	if err := params.Config.Headers.Secrets(params.FileSystem); err != nil {
		return nil, se.Prefix("tracer", err)
	}

	client := otlp.NewClient(otlp.WithEndpointURL(params.Config.URL), otlp.WithHeaders(params.Config.Headers))
	exporter := otlptrace.NewUnstarted(client)

	attrs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(params.Name.String()),
		semconv.ServiceVersion(params.Version.String()),
		semconv.DeploymentEnvironmentName(params.Environment.String()),
	)

	provider := sdktrace.NewTracerProvider(sdktrace.WithResource(attrs), sdktrace.WithBatcher(exporter))

	otel.SetTracerProvider(provider)
	otel.SetErrorHandler(&errorHandler{logger: params.Logger})

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return se.Prefix("tracer", exporter.Start(ctx))
		},
		OnStop: func(ctx context.Context) error {
			_ = provider.Shutdown(ctx)
			_ = exporter.Shutdown(ctx)

			return nil
		},
	})

	return &Tracer{provider.Tracer(params.Name.String())}, nil
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

type errorHandler struct {
	logger *logger.Logger
}

func (e *errorHandler) Handle(err error) {
	e.logger.Error("tracer: global error", zap.Error(err))
}

package tracer

import (
	"context"

	"github.com/alexfalkowski/go-service/env"
	se "github.com/alexfalkowski/go-service/errors"
	"go.opentelemetry.io/otel"
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

// NewTracer for tracer.
func NewTracer(lc fx.Lifecycle, env env.Environment, ver env.Version, name env.Name, cfg *Config, logger *zap.Logger) (trace.Tracer, error) {
	if !IsEnabled(cfg) {
		return noop.Tracer{}, nil
	}

	if err := cfg.Headers.Secrets(); err != nil {
		return nil, se.Prefix("tracer", err)
	}

	client := otlp.NewClient(otlp.WithEndpointURL(cfg.URL), otlp.WithHeaders(cfg.Headers))
	exporter := otlptrace.NewUnstarted(client)

	attrs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(string(name)),
		semconv.ServiceVersion(string(ver)),
		semconv.DeploymentEnvironmentName(string(env)),
	)

	p := sdktrace.NewTracerProvider(sdktrace.WithResource(attrs), sdktrace.WithBatcher(exporter))

	otel.SetTracerProvider(p)
	otel.SetErrorHandler(&errorHandler{logger: logger})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return se.Prefix("tracer", exporter.Start(ctx))
		},
		OnStop: func(ctx context.Context) error {
			_ = p.Shutdown(ctx)
			_ = exporter.Shutdown(ctx)

			return nil
		},
	})

	return p.Tracer(string(name)), nil
}

type errorHandler struct {
	logger *zap.Logger
}

func (e *errorHandler) Handle(err error) {
	e.logger.Error("tracer: global error", zap.Error(err))
}

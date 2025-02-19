package tracer

import (
	"context"

	"github.com/alexfalkowski/go-service/env"
	se "github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/telemetry/logger"
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
func NewTracer(lc fx.Lifecycle, env env.Environment, ver env.Version, name env.Name, fs os.FileSystem, cfg *Config, logger *logger.Logger) (trace.Tracer, error) {
	if !IsEnabled(cfg) {
		return noop.Tracer{}, nil
	}

	if err := cfg.Headers.Secrets(fs); err != nil {
		return nil, se.Prefix("tracer", err)
	}

	client := otlp.NewClient(otlp.WithEndpointURL(cfg.URL), otlp.WithHeaders(cfg.Headers))
	exporter := otlptrace.NewUnstarted(client)

	attrs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(name.String()),
		semconv.ServiceVersion(ver.String()),
		semconv.DeploymentEnvironmentName(env.String()),
	)

	provider := sdktrace.NewTracerProvider(sdktrace.WithResource(attrs), sdktrace.WithBatcher(exporter))

	otel.SetTracerProvider(provider)
	otel.SetErrorHandler(&errorHandler{logger: logger})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return se.Prefix("tracer", exporter.Start(ctx))
		},
		OnStop: func(ctx context.Context) error {
			_ = provider.Shutdown(ctx)
			_ = exporter.Shutdown(ctx)

			return nil
		},
	})

	return provider.Tracer(name.String()), nil
}

type errorHandler struct {
	logger *logger.Logger
}

func (e *errorHandler) Handle(err error) {
	e.logger.Error("tracer: global error", zap.Error(err))
}

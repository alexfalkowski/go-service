package tracer

import (
	"context"
	"errors"

	"github.com/alexfalkowski/go-service/env"
	se "github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/net"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	otlp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
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
func NewTracer(lc fx.Lifecycle, env env.Environment, ver version.Version, cfg *Config, logger *zap.Logger) trace.Tracer {
	if !IsEnabled(cfg) {
		return noop.Tracer{}
	}

	opts := []otlp.Option{otlp.WithEndpointURL(cfg.Host)}
	if cfg.HasKey() {
		opts = append(opts, otlp.WithHeaders(map[string]string{"Authorization": "Basic " + cfg.GetKey()}))
	}

	return newTracer(lc, env, ver, logger, opts)
}

func newTracer(lc fx.Lifecycle, env env.Environment, ver version.Version, logger *zap.Logger, opts []otlp.Option) trace.Tracer {
	client := otlp.NewClient(opts...)
	exporter := otlptrace.NewUnstarted(client)

	name := os.ExecutableName()
	attrs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(name),
		semconv.ServiceVersion(string(ver)),
		semconv.DeploymentEnvironment(string(env)),
	)

	p := sdktrace.NewTracerProvider(sdktrace.WithResource(attrs), sdktrace.WithBatcher(exporter))

	otel.SetTracerProvider(p)
	otel.SetErrorHandler(&errorHandler{logger: logger})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return se.Prefix("start tracer", exporter.Start(ctx))
		},
		OnStop: func(ctx context.Context) error {
			return se.Prefix("stop tracer", errors.Join(p.Shutdown(ctx), exporter.Shutdown(ctx)))
		},
	})

	return p.Tracer(name)
}

type errorHandler struct {
	logger *zap.Logger
}

func (e *errorHandler) Handle(err error) {
	if net.IsConnectionRefused(err) {
		return
	}

	e.logger.Error("trace: global error", zap.Error(err))
}

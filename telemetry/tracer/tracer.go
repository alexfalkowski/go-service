package tracer

import (
	"context"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// Register for tracer.
func Register() {
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

// NewNoopTracer for tracer.
func NewNoopTracer(name string) trace.Tracer {
	return trace.NewNoopTracerProvider().Tracer(name)
}

// NewTracer for tracer.
func NewTracer(lc fx.Lifecycle, name string, env env.Environment, ver version.Version, cfg *Config) (trace.Tracer, error) {
	if cfg.Host == "" {
		return NewNoopTracer(name), nil
	}

	opts := []otlptracehttp.Option{otlptracehttp.WithEndpoint(cfg.Host)}

	if !cfg.Secure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	client := otlptracehttp.NewClient(opts...)

	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return nil, err
	}

	attrs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(name),
		semconv.ServiceVersion(string(ver)),
		semconv.DeploymentEnvironment(string(env)),
		attribute.String("name", os.ExecutableName()),
	)

	tracerOpts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(attrs),
	}

	if env.IsDevelopment() {
		tracerOpts = append(tracerOpts, sdktrace.WithSyncer(exporter))
	} else {
		tracerOpts = append(tracerOpts, sdktrace.WithBatcher(exporter))
	}

	p := sdktrace.NewTracerProvider(tracerOpts...)

	otel.SetTracerProvider(p)
	otel.SetErrorHandler(&errorHandler{})

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return p.Shutdown(ctx)
		},
	})

	return p.Tracer(name), nil
}

type errorHandler struct{}

func (*errorHandler) Handle(_ error) {
}

package tracer

import (
	"context"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
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

// NewNoopTracer for tracer.
func NewNoopTracer() trace.Tracer {
	return noop.Tracer{}
}

// NewTracer for tracer.
func NewTracer(lc fx.Lifecycle, env env.Environment, ver version.Version, cfg *Config, logger *zap.Logger) (trace.Tracer, error) {
	if !IsEnabled(cfg) {
		return NewNoopTracer(), nil
	}

	opts := []otlptracehttp.Option{}

	if cfg.IsBaselime() {
		opts = append(opts, otlptracehttp.WithEndpointURL("https://otel.baselime.io"), otlptracehttp.WithHeaders(map[string]string{"x-api-key": cfg.Key}))
	} else {
		opts = append(opts, otlptracehttp.WithEndpointURL(cfg.Host))
	}

	return newTracer(context.Background(), lc, env, ver, logger, opts)
}

func newTracer(ctx context.Context, lc fx.Lifecycle, env env.Environment, ver version.Version, logger *zap.Logger, opts []otlptracehttp.Option) (trace.Tracer, error) {
	client := otlptracehttp.NewClient(opts...)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, err
	}

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
		OnStop: func(ctx context.Context) error {
			return p.Shutdown(ctx)
		},
	})

	return p.Tracer(name), nil
}

type errorHandler struct {
	logger *zap.Logger
}

func (e *errorHandler) Handle(err error) {
	e.logger.Error("trace: global error", zap.Error(err))
}

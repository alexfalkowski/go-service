package tracer

import (
	"context"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	"github.com/honeycombio/otel-config-go/otelconfig"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/fx"
)

// Register for tracer.
func Register() {
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

// NewNoopTracer for tracer.
func NewNoopTracer(name string) trace.Tracer {
	return noop.NewTracerProvider().Tracer(name)
}

// NewTracer for tracer.
func NewTracer(ctx context.Context, lc fx.Lifecycle, name string, env env.Environment, ver version.Version, cfg *Config) (trace.Tracer, error) {
	if cfg == nil || !cfg.Enabled {
		return NewNoopTracer(name), nil
	}

	if cfg.IsBaselime() {
		return baselimeTracer(lc, name, cfg)
	}

	return defaultTracer(ctx, lc, name, env, ver, cfg)
}

func baselimeTracer(lc fx.Lifecycle, name string, cfg *Config) (trace.Tracer, error) {
	e := cfg.Host
	if e == "" {
		e = "https://otel.baselime.io"
	}

	sh, err := otelconfig.ConfigureOpenTelemetry(
		otelconfig.WithLogger(&logger{}),
		otelconfig.WithExporterProtocol(otelconfig.ProtocolHTTPProto),
		otelconfig.WithExporterEndpoint(e),
		otelconfig.WithServiceName(name),
		otelconfig.WithHeaders(map[string]string{"x-api-key": cfg.Key}),
	)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			sh()

			return nil
		},
	})

	return otel.Tracer(name), nil
}

func defaultTracer(ctx context.Context, lc fx.Lifecycle, name string, env env.Environment, ver version.Version, cfg *Config) (trace.Tracer, error) {
	opts := []otlptracehttp.Option{otlptracehttp.WithEndpoint(cfg.Host)}

	if !cfg.Secure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	client := otlptracehttp.NewClient(opts...)

	exporter, err := otlptrace.New(ctx, client)
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

	p := sdktrace.NewTracerProvider(sdktrace.WithResource(attrs), sdktrace.WithBatcher(exporter))

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

type logger struct{}

func (l *logger) Fatalf(_ string, _ ...interface{}) {
}

func (l *logger) Debugf(_ string, _ ...interface{}) {
}

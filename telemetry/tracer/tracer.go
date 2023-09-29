package tracer

import (
	"context"

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

// Params for tracer.
type Params struct {
	Lifecycle fx.Lifecycle
	Name      string
	Version   version.Version
	Config    *Config
}

// Register for tracer.
func Register() {
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

// NewNoopTracer for tracer.
func NewNoopTracer(name string) trace.Tracer {
	return trace.NewNoopTracerProvider().Tracer(name)
}

// NewTracer for tracer.
func NewTracer(params Params) (trace.Tracer, error) {
	if params.Config.Host == "" {
		return NewNoopTracer(params.Name), nil
	}

	opts := []otlptracehttp.Option{otlptracehttp.WithEndpoint(params.Config.Host)}

	if !params.Config.Secure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	client := otlptracehttp.NewClient(opts...)

	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return nil, err
	}

	attrs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(params.Name),
		semconv.ServiceVersion(string(params.Version)),
		attribute.String("name", os.ExecutableName()),
	)

	p := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(attrs),
	)

	otel.SetTracerProvider(p)
	otel.SetErrorHandler(&errorHandler{})

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return p.Shutdown(ctx)
		},
	})

	return p.Tracer(params.Name), nil
}

type errorHandler struct{}

func (*errorHandler) Handle(_ error) {
}

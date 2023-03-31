package otel

import (
	"context"
	"net"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// TracerParams for otel.
type TracerParams struct {
	Lifecycle fx.Lifecycle
	Name      string
	Version   version.Version
	Config    *Config
}

// Register for otel.
func Register() {
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

// NewNoopTracer for otel.
func NewNoopTracer(name string) trace.Tracer {
	return trace.NewNoopTracerProvider().Tracer(name)
}

// NewTracer for otel.
func NewTracer(params TracerParams) (trace.Tracer, error) {
	if params.Config.IsJaeger() {
		host, port, err := net.SplitHostPort(params.Config.Host)
		if err != nil {
			return nil, err
		}

		exp, err := jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost(host), jaeger.WithAgentPort(port)))
		if err != nil {
			return nil, err
		}

		p := tracesdk.NewTracerProvider(
			tracesdk.WithBatcher(exp),
			tracesdk.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(params.Name),
				semconv.ServiceVersion(string(params.Version)),
				attribute.String("name", os.ExecutableName()),
			)),
		)

		params.Lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return p.Shutdown(ctx)
			},
		})

		return p.Tracer(params.Name), nil
	}

	return NewNoopTracer(params.Name), nil
}

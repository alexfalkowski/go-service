package tracer

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	otlp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// TracerParams for tracer.
type TracerParams struct {
	di.In
	Lifecycle   di.Lifecycle
	Config      *Config
	ID          env.ID
	Name        env.Name
	Version     env.Version
	Environment env.Environment
}

// NewTracer for tracer.
func NewTracer(params TracerParams) *Tracer {
	if !params.Config.IsEnabled() {
		return nil
	}

	client := otlp.NewClient(otlp.WithEndpointURL(params.Config.URL), otlp.WithHeaders(params.Config.Headers))
	exporter := otlptrace.NewUnstarted(client)
	attrs := resource.NewWithAttributes(
		attributes.SchemaURL,
		attributes.HostID(params.ID.String()),
		attributes.ServiceName(params.Name.String()),
		attributes.ServiceVersion(params.Version.String()),
		attributes.DeploymentEnvironmentName(params.Environment.String()),
	)

	provider := sdk.NewTracerProvider(sdk.WithResource(attrs), sdk.WithBatcher(exporter))
	otel.SetTracerProvider(provider)

	params.Lifecycle.Append(di.Hook{
		OnStart: func(ctx context.Context) error {
			return prefix(exporter.Start(ctx))
		},
		OnStop: func(ctx context.Context) error {
			// Do not return error as this will stop all others.
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

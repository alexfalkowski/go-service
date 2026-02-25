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
)

// TracerParams declares the dependencies required by Register.
//
// It is intended for Fx/Dig injection and includes service identity fields used to
// populate OpenTelemetry resource attributes.
type TracerParams struct {
	di.In

	// Lifecycle is used to start and stop the exporter/provider with the application.
	Lifecycle di.Lifecycle

	// Config enables tracing when non-nil and supplies exporter settings.
	Config *Config

	// ID is the host identifier used for the resource's host.id attribute.
	ID env.ID

	// Name is the service name used for the resource's service.name attribute.
	Name env.Name

	// Version is the service version used for the resource's service.version attribute.
	Version env.Version

	// Environment is the deployment environment name used for the resource's
	// deployment.environment.name attribute.
	Environment env.Environment
}

// Register configures and installs a global OpenTelemetry TracerProvider.
//
// When tracing is enabled (`params.Config != nil`), Register:
//
//  1. Creates an OTLP/HTTP trace exporter using `Config.URL` and `Config.Headers`.
//  2. Creates an OpenTelemetry resource describing the running service.
//  3. Installs a `sdk.TracerProvider` via `otel.SetTracerProvider`.
//  4. Appends lifecycle hooks to start the exporter on application start and to
//     shut down the provider/exporter on application stop.
//
// Shutdown errors are intentionally ignored to avoid blocking other stop hooks.
func Register(params TracerParams) {
	if !params.Config.IsEnabled() {
		return
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
}

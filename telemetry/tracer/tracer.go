package tracer

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/internal/otlp"
	"github.com/alexfalkowski/go-sync"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

var (
	enabled      sync.Bool
	noopProvider trace.TracerProvider = noop.NewTracerProvider()
)

// ReadOnlySpan is an alias for [go.opentelemetry.io/otel/sdk/trace.ReadOnlySpan].
type ReadOnlySpan = sdk.ReadOnlySpan

// SpanExporter is an alias for [go.opentelemetry.io/otel/sdk/trace.SpanExporter].
type SpanExporter = sdk.SpanExporter

// Provider is an alias for [trace.TracerProvider].
type Provider = trace.TracerProvider

// SDKProvider is an alias for [go.opentelemetry.io/otel/sdk/trace.TracerProvider].
type SDKProvider = sdk.TracerProvider

// ProviderOption is an alias for [go.opentelemetry.io/otel/sdk/trace.TracerProviderOption].
type ProviderOption = sdk.TracerProviderOption

// NewProvider constructs an OpenTelemetry SDK tracer provider.
func NewProvider(opts ...ProviderOption) *SDKProvider {
	return sdk.NewTracerProvider(opts...)
}

// WithSyncer configures a tracer provider to export spans synchronously.
func WithSyncer(exporter SpanExporter) ProviderOption {
	return sdk.WithSyncer(exporter)
}

// NewNoopProvider constructs a no-op tracer provider.
func NewNoopProvider() Provider {
	return noop.NewTracerProvider()
}

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

// Register configures and installs a global OpenTelemetry tracer provider.
//
// When tracing is configured with kind "otlp", Register:
//
//  1. Creates an OTLP/HTTP trace exporter using [Config.URL] and [Config.Headers].
//  2. Creates an OpenTelemetry resource describing the running service.
//  3. Installs a [go.opentelemetry.io/otel/sdk/trace.TracerProvider] globally.
//  4. Appends lifecycle hooks to start the exporter on application start and to
//     shut down the provider/exporter on application stop.
//
// If Config is nil or Kind is empty, Register installs the noop provider. Unknown
// non-empty kinds return ErrNotFound.
//
// Shutdown errors are intentionally ignored to avoid blocking other stop hooks.
func Register(params TracerParams) error {
	if !params.Config.IsEnabled() {
		setProvider(noopProvider, false)
		return nil
	}

	switch params.Config.Kind {
	case "otlp":
		if err := otlp.ValidateEndpoint(params.Config.URL, params.Config.Headers); err != nil {
			return prefix(err)
		}

		opts := []otlptracehttp.Option{otlptracehttp.WithHeaders(params.Config.Headers)}
		if !strings.IsEmpty(params.Config.URL) {
			opts = append(opts, otlptracehttp.WithEndpointURL(params.Config.URL))
		}

		client := otlptracehttp.NewClient(opts...)
		exporter := otlptrace.NewUnstarted(client)
		attrs := resource.NewWithAttributes(
			attributes.SchemaURL,
			attributes.HostID(params.ID.String()),
			attributes.ServiceName(params.Name.String()),
			attributes.ServiceVersion(params.Version.String()),
			attributes.DeploymentEnvironmentName(params.Environment.String()),
		)

		provider := sdk.NewTracerProvider(sdk.WithResource(attrs), sdk.WithBatcher(exporter))
		setProvider(provider, true)

		params.Lifecycle.Append(di.Hook{
			OnStart: func(ctx context.Context) error {
				return prefix(exporter.Start(ctx))
			},
			OnStop: func(ctx context.Context) error {
				// Do not return error as this will stop all others.
				_ = provider.Shutdown(ctx)
				_ = exporter.Shutdown(ctx)
				setProvider(noopProvider, false)

				return nil
			},
		})

		return nil
	default:
		return ErrNotFound
	}
}

// IsEnabled reports whether this package has registered tracing as enabled.
func IsEnabled() bool {
	return enabled.Load()
}

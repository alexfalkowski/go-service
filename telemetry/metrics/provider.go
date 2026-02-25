package metrics

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

// MeterProviderParams declares the dependencies required by NewMeterProvider.
//
// It is intended for Fx/Dig injection and includes service identity fields used to
// populate OpenTelemetry resource attributes.
type MeterProviderParams struct {
	di.In

	// Lifecycle is used to start runtime/host instrumentation and to shut down the
	// meter provider with the application.
	Lifecycle di.Lifecycle

	// Config enables metrics when non-nil and supplies exporter settings.
	Config *Config

	// Reader is the SDK reader/exporter that NewMeterProvider will attach to the provider.
	// A nil reader disables metrics even when Config is set.
	Reader sdk.Reader

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

// NewMeterProvider constructs and installs a global OpenTelemetry MeterProvider.
//
// When metrics are enabled (`params.Config != nil`) and a non-nil Reader is provided,
// NewMeterProvider:
//
//  1. Creates an OpenTelemetry resource describing the running service.
//  2. Constructs an SDK `*metric.MeterProvider` with the provided reader.
//  3. Installs it globally via `otel.SetMeterProvider`.
//  4. Registers lifecycle hooks:
//     - OnStart: starts runtime and host instrumentation using this provider
//     - OnStop: shuts down the provider
//
// Provider shutdown errors are intentionally ignored to avoid blocking other stop hooks.
//
// If metrics are disabled or Reader is nil, it returns nil.
func NewMeterProvider(params MeterProviderParams) MeterProvider {
	if !params.Config.IsEnabled() || params.Reader == nil {
		return nil
	}

	reader := params.Reader
	attrs := resource.NewWithAttributes(
		attributes.SchemaURL,
		attributes.HostID(params.ID.String()),
		attributes.ServiceName(params.Name.String()),
		attributes.ServiceVersion(params.Version.String()),
		attributes.DeploymentEnvironmentName(params.Environment.String()),
	)
	provider := sdk.NewMeterProvider(sdk.WithReader(reader), sdk.WithResource(attrs))
	otel.SetMeterProvider(provider)

	params.Lifecycle.Append(di.Hook{
		OnStart: func(_ context.Context) error {
			err := errors.Join(runtime.Start(runtime.WithMeterProvider(provider)), host.Start(host.WithMeterProvider(provider)))

			return prefix(err)
		},
		OnStop: func(ctx context.Context) error {
			// Do not return error as this will stop all others.
			_ = provider.Shutdown(ctx)

			return nil
		},
	})

	return provider
}

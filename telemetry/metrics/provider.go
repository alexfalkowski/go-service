package metrics

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-sync"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	sdk "go.opentelemetry.io/otel/sdk/metric"
)

var (
	enabled      sync.Bool
	noopProvider MeterProvider = NewNoopMeterProvider()
)

// MeterProviderParams declares the dependencies required by NewMeterProvider.
//
// It is intended for Fx/Dig injection and includes service identity fields used to
// populate OpenTelemetry resource attributes.
type MeterProviderParams struct {
	di.In

	// Lifecycle is used to start runtime instrumentation and to shut down the meter
	// provider with the application.
	Lifecycle di.Lifecycle

	// Reader is the SDK reader/exporter that NewMeterProvider will attach to the provider.
	// A nil reader disables metrics even when Config is set.
	Reader sdk.Reader

	// Config enables metrics when non-nil and supplies exporter settings.
	Config *Config

	// Attributes are optional OpenTelemetry resource attributes attached to metrics.
	Attributes attributes.Map `optional:"true"`

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
//  2. Constructs an SDK *[go.opentelemetry.io/otel/sdk/metric.MeterProvider] with the provided reader.
//  3. Installs it globally via [go.opentelemetry.io/otel.SetMeterProvider].
//  4. Registers lifecycle hooks:
//     - OnStart: starts runtime instrumentation using this provider
//     - OnStop: shuts down the provider
//
// Provider shutdown errors are intentionally ignored to avoid blocking other stop hooks.
//
// If metrics are disabled or Reader is nil, it installs and returns the package noop provider.
func NewMeterProvider(params MeterProviderParams) MeterProvider {
	if !params.Config.IsEnabled() || params.Reader == nil {
		setMeterProvider(noopProvider, false)
		return noopProvider
	}

	attrs := attributes.NewResource(
		params.Attributes,
		params.ID.String(),
		params.Name.String(),
		params.Version.String(),
		params.Environment.String(),
	)
	options := []sdk.Option{sdk.WithReader(params.Reader), sdk.WithResource(attrs)}
	if views := configViews(params.Config); len(views) > 0 {
		options = append(options, sdk.WithView(views...))
	}

	provider := sdk.NewMeterProvider(options...)
	setMeterProvider(provider, true)

	params.Lifecycle.Append(di.Hook{
		OnStart: func(_ context.Context) error {
			// Re-add host metrics when https://github.com/shirou/gopsutil/issues/2115 is fixed.
			err := runtime.Start(runtime.WithMeterProvider(provider))

			return prefix(err)
		},
		OnStop: func(ctx context.Context) error {
			// Do not return error as this will stop all others.
			_ = provider.Shutdown(ctx)
			setMeterProvider(noopProvider, false)

			return nil
		},
	})

	return provider
}

// configViews builds explicit-bucket-histogram views from the configured
// instrument-name-to-boundaries map. It returns nil when no views are
// configured, leaving the OpenTelemetry SDK default boundaries in place.
func configViews(cfg *Config) []sdk.View {
	if cfg == nil || len(cfg.Views) == 0 {
		return nil
	}

	views := make([]sdk.View, 0, len(cfg.Views))
	for name, boundaries := range cfg.Views {
		views = append(views, sdk.NewView(
			sdk.Instrument{Name: name},
			sdk.Stream{Aggregation: sdk.AggregationExplicitBucketHistogram{Boundaries: boundaries}},
		))
	}

	return views
}

// IsEnabled reports whether this package has registered metrics as enabled.
func IsEnabled() bool {
	return enabled.Load()
}

// NewManualReader constructs an OpenTelemetry SDK manual metric reader.
func NewManualReader() sdk.Reader {
	return sdk.NewManualReader()
}

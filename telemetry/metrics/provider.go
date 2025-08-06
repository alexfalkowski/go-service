package metrics

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

// MeterProviderParams for metrics.
type MeterProviderParams struct {
	di.In
	Lifecycle   di.Lifecycle
	Config      *Config
	Reader      sdk.Reader
	ID          env.ID
	Name        env.Name
	Version     env.Version
	Environment env.Environment
}

// NewMeterProvider for metrics.
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

	params.Lifecycle.Append(di.Hook{
		OnStart: func(_ context.Context) error {
			err := errors.Join(runtime.Start(runtime.WithMeterProvider(provider)), host.Start(host.WithMeterProvider(provider)))

			return prefix(err)
		},
		OnStop: func(ctx context.Context) error {
			_ = provider.Shutdown(ctx)
			_ = reader.Shutdown(ctx)

			return nil
		},
	})

	return provider
}

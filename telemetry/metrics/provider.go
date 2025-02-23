package metrics

import (
	"context"
	"errors"

	"github.com/alexfalkowski/go-service/env"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	om "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	sm "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.uber.org/fx"
)

// MeterProviderParams for metrics.
type MeterProviderParams struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Config      *Config
	Reader      sm.Reader
	Environment env.Environment
	Version     env.Version
	Name        env.Name
}

// NewMeterProvider for metrics.
func NewMeterProvider(params MeterProviderParams) om.MeterProvider {
	if !IsEnabled(params.Config) || params.Reader == nil {
		return &noop.MeterProvider{}
	}

	name := params.Name.String()
	attrs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(name),
		semconv.ServiceVersion(params.Version.String()),
		semconv.DeploymentEnvironmentName(params.Environment.String()),
	)
	provider := sm.NewMeterProvider(sm.WithReader(params.Reader), sm.WithResource(attrs))

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			err := errors.Join(runtime.Start(runtime.WithMeterProvider(provider)), host.Start(host.WithMeterProvider(provider)))

			return prefix(err)
		},
		OnStop: func(ctx context.Context) error {
			_ = provider.Shutdown(ctx)

			return nil
		},
	})

	return provider
}

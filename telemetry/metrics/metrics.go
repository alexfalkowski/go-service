package metrics

import (
	"context"
	"errors"

	"github.com/alexfalkowski/go-service/env"
	se "github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/os"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	otlp "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	om "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	sm "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.uber.org/fx"
)

// NewReader for metrics. A nil reader means disabled.
func NewReader(fs os.FileSystem, cfg *Config) (sm.Reader, error) {
	switch {
	case !IsEnabled(cfg):
		return nil, nil
	case cfg.IsOTLP():
		if err := cfg.Headers.Secrets(fs); err != nil {
			return nil, se.Prefix("metrics", err)
		}

		r, err := otlp.New(context.Background(), otlp.WithEndpointURL(cfg.URL), otlp.WithHeaders(cfg.Headers))

		return sm.NewPeriodicReader(r), se.Prefix("metrics", err)
	case cfg.IsPrometheus():
		e, err := prometheus.New()

		return e, se.Prefix("metrics", err)
	default:
		return nil, nil
	}
}

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

			return se.Prefix("metrics", err)
		},
		OnStop: func(ctx context.Context) error {
			_ = provider.Shutdown(ctx)

			return nil
		},
	})

	return provider
}

// MeterParams for metrics.
type MeterParams struct {
	fx.In

	Config   *Config
	Provider om.MeterProvider
	Name     env.Name
}

// NewMeter for metrics.
func NewMeter(provider om.MeterProvider, name env.Name) *Meter {
	return &Meter{provider.Meter(name.String())}
}

// Meter using otel.
type Meter struct {
	om.Meter
}

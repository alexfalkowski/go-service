package metrics

import (
	"context"
	"errors"

	"github.com/alexfalkowski/go-service/env"
	se "github.com/alexfalkowski/go-service/errors"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	otlp "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	om "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	sm "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
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
	if !IsEnabled(params.Config) {
		return &noop.MeterProvider{}
	}

	name := string(params.Name)
	attrs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(name),
		semconv.ServiceVersion(string(params.Version)),
		semconv.DeploymentEnvironment(string(params.Environment)),
	)
	provider := sm.NewMeterProvider(sm.WithReader(params.Reader), sm.WithResource(attrs))

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			err := errors.Join(runtime.Start(runtime.WithMeterProvider(provider)), host.Start(host.WithMeterProvider(provider)))

			return se.Prefix("new meter", err)
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
func NewMeter(provider om.MeterProvider, name env.Name) om.Meter {
	return provider.Meter(string(name))
}

// NewReader for metrics. A nil reader means disabled.
//
//nolint:nilnil
func NewReader(cfg *Config) (sm.Reader, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	if cfg.IsOTLP() {
		opts := []otlp.Option{otlp.WithEndpointURL(cfg.URL)}

		if cfg.HasKey() {
			k, err := cfg.GetKey()
			if err != nil {
				return nil, err
			}

			opts = append(opts, otlp.WithHeaders(map[string]string{"Authorization": "Basic " + k}))
		}

		r, err := otlp.New(context.Background(), opts...)

		return sm.NewPeriodicReader(r), se.Prefix("new otlp", err)
	}

	e, err := prometheus.New()

	return e, se.Prefix("new prometheus", err)
}

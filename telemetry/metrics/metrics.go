package metrics

import (
	"context"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	otlp "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	m "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.uber.org/fx"
)

// NewNoopMeter for metrics.
func NewNoopMeter() m.Meter {
	return noop.Meter{}
}

// MeterParams for metrics.
type MeterParams struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Environment env.Environment
	Version     version.Version
	Config      *Config
	Reader      metric.Reader
}

// NewMeter for metrics.
func NewMeter(params MeterParams) (m.Meter, error) {
	if !IsEnabled(params.Config) {
		return NewNoopMeter(), nil
	}

	name := os.ExecutableName()
	attrs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(name),
		semconv.ServiceVersion(string(params.Version)),
		semconv.DeploymentEnvironment(string(params.Environment)),
	)

	provider := metric.NewMeterProvider(metric.WithReader(params.Reader), metric.WithResource(attrs))
	meter := provider.Meter(name)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return provider.Shutdown(ctx)
		},
	})

	return meter, nil
}

// NewReader for metrics.
func NewReader(cfg *Config) (metric.Reader, error) {
	if !IsEnabled(cfg) {
		return prometheus.New()
	}

	if cfg.IsOTLP() {
		r, err := otlp.New(context.Background(), otlp.WithEndpointURL(cfg.Host))

		return metric.NewPeriodicReader(r), err
	}

	return prometheus.New()
}

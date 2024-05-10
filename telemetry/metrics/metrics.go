package metrics

import (
	"context"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/net"
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
	Reader      metric.Reader
	Config      *Config
	Environment env.Environment
	Version     version.Version
}

// NewMeter for metrics.
func NewMeter(params MeterParams) m.Meter {
	if !IsEnabled(params.Config) {
		return NewNoopMeter()
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
			err := provider.Shutdown(ctx)
			if net.IsConnectionRefused(err) {
				return nil
			}

			return errors.Prefix("stop metrics", err)
		},
	})

	return meter
}

// NewReader for metrics.
func NewReader(cfg *Config) (metric.Reader, error) {
	if !IsEnabled(cfg) {
		return prom()
	}

	if cfg.IsOTLP() {
		r, err := otlp.New(context.Background(), otlp.WithEndpointURL(cfg.Host))

		return metric.NewPeriodicReader(r), errors.Prefix("new otlp", err)
	}

	return prom()
}

func prom() (*prometheus.Exporter, error) {
	e, err := prometheus.New()

	return e, errors.Prefix("new prometheus", err)
}

package metrics

import (
	"context"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/errors"
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
	Version     env.Version
	Name        env.Name
}

// NewMeter for metrics.
func NewMeter(params MeterParams) m.Meter {
	if !IsEnabled(params.Config) {
		return NewNoopMeter()
	}

	name := string(params.Name)
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
			_ = provider.Shutdown(ctx)

			return nil
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
		opts := []otlp.Option{otlp.WithEndpointURL(cfg.Host)}

		if cfg.HasKey() {
			k, err := cfg.GetKey()
			if err != nil {
				return nil, err
			}

			opts = append(opts, otlp.WithHeaders(map[string]string{"Authorization": "Basic " + k}))
		}

		r, err := otlp.New(context.Background(), opts...)

		return metric.NewPeriodicReader(r), errors.Prefix("new otlp", err)
	}

	return prom()
}

func prom() (*prometheus.Exporter, error) {
	e, err := prometheus.New()

	return e, errors.Prefix("new prometheus", err)
}

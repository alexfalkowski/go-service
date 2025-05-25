package test

import (
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
)

// NewOTLPMeter for test.
func NewOTLPMeter(lc fx.Lifecycle) *metrics.Meter {
	return NewMeter(lc, NewOTLPMetricsConfig())
}

// NewPrometheusMeter for test.
func NewPrometheusMeter(lc fx.Lifecycle) *metrics.Meter {
	return NewMeter(lc, NewPrometheusMetricsConfig())
}

// NewMeter for test.
func NewMeter(lc fx.Lifecycle, c *metrics.Config) *metrics.Meter {
	return metrics.NewMeter(NewMeterProvider(lc, c), Name)
}

// NewOTLPMeterProvider for test.
func NewOTLPMeterProvider(lc fx.Lifecycle) metric.MeterProvider {
	return NewMeterProvider(lc, NewOTLPMetricsConfig())
}

// NewPrometheusMeterProvider for test.
func NewPrometheusMeterProvider(lc fx.Lifecycle) metric.MeterProvider {
	return NewMeterProvider(lc, NewPrometheusMetricsConfig())
}

// NewMeterProvider for test.
func NewMeterProvider(lc fx.Lifecycle, config *metrics.Config) metric.MeterProvider {
	r, err := metrics.NewReader(config)
	runtime.Must(err)

	params := metrics.MeterProviderParams{
		Lifecycle:   lc,
		Config:      config,
		Reader:      r,
		Environment: Environment,
		Version:     Version,
		Name:        Name,
	}

	return metrics.NewMeterProvider(params)
}

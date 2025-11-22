package test

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"go.opentelemetry.io/otel/metric"
)

// NewOTLPMeter for test.
func NewOTLPMeter(lc di.Lifecycle) metrics.Meter {
	return NewMeter(lc, NewOTLPMetricsConfig())
}

// NewPrometheusMeter for test.
func NewPrometheusMeter(lc di.Lifecycle) metrics.Meter {
	return NewMeter(lc, NewPrometheusMetricsConfig())
}

// NewMeter for test.
func NewMeter(lc di.Lifecycle, c *metrics.Config) metrics.Meter {
	return metrics.NewMeter(Name, Version, NewMeterProvider(lc, c))
}

// NewOTLPMeterProvider for test.
func NewOTLPMeterProvider(lc di.Lifecycle) metric.MeterProvider {
	return NewMeterProvider(lc, NewOTLPMetricsConfig())
}

// NewPrometheusMeterProvider for test.
func NewPrometheusMeterProvider(lc di.Lifecycle) metric.MeterProvider {
	return NewMeterProvider(lc, NewPrometheusMetricsConfig())
}

// NewMeterProvider for test.
func NewMeterProvider(lc di.Lifecycle, config *metrics.Config) metric.MeterProvider {
	r, err := metrics.NewReader(lc, Name, config)
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

package test

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"go.opentelemetry.io/otel/metric"
)

// NewOTLPMeter returns a meter backed by the shared OTLP metrics config.
func NewOTLPMeter(lc di.Lifecycle) metrics.Meter {
	return NewMeter(lc, NewOTLPMetricsConfig())
}

// NewPrometheusMeter returns a meter backed by the shared Prometheus metrics config.
func NewPrometheusMeter(lc di.Lifecycle) metrics.Meter {
	return NewMeter(lc, NewPrometheusMetricsConfig())
}

// NewMeter returns a repository meter scoped to the shared test name and version.
func NewMeter(lc di.Lifecycle, c *metrics.Config) metrics.Meter {
	return metrics.NewMeter(Name, Version, NewMeterProvider(lc, c))
}

// NewOTLPMeterProvider returns a meter provider backed by the shared OTLP metrics config.
func NewOTLPMeterProvider(lc di.Lifecycle) metric.MeterProvider {
	return NewMeterProvider(lc, NewOTLPMetricsConfig())
}

// NewPrometheusMeterProvider returns a meter provider backed by the shared Prometheus metrics config.
func NewPrometheusMeterProvider(lc di.Lifecycle) metric.MeterProvider {
	return NewMeterProvider(lc, NewPrometheusMetricsConfig())
}

// NewMeterProvider creates a meter provider with a reader registered on the supplied lifecycle.
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

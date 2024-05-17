package test

import (
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
)

// NewOTLPMeter for test.
func NewOTLPMeter(lc fx.Lifecycle) metric.Meter {
	return NewMeter(lc, NewOTLPMetricsConfig())
}

// NewPrometheusMeter for test.
func NewPrometheusMeter(lc fx.Lifecycle) metric.Meter {
	return NewMeter(lc, NewPrometheusMetricsConfig())
}

// NewMeter for test.
func NewMeter(lc fx.Lifecycle, c *metrics.Config) metric.Meter {
	r, err := metrics.NewReader(c)
	runtime.Must(err)

	p := metrics.MeterParams{
		Lifecycle:   lc,
		Environment: Environment,
		Version:     Version,
		Config:      c,
		Reader:      r,
	}

	return metrics.NewMeter(p)
}

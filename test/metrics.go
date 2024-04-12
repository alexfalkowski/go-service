package test

import (
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
)

// NewMeter for test.
func NewMeter(lc fx.Lifecycle) metric.Meter {
	m, err := metrics.NewMeter(lc, "test", Environment, Version, NewPrometheusMetricsConfig())
	if err != nil {
		panic(err)
	}

	return m
}

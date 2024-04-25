package test

import (
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
)

// NewMeter for test.
func NewMeter(lc fx.Lifecycle) metric.Meter {
	m, err := metrics.NewMeter(lc, Environment, Version, NewPrometheusMetricsConfig())
	runtime.Must(err)

	return m
}

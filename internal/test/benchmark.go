package test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

// ResetTelemetry disables global telemetry state for benchmarks and tests.
func ResetTelemetry(tb testing.TB) {
	tb.Helper()

	require.NoError(tb, tracer.Register(tracer.TracerParams{Lifecycle: fxtest.NewLifecycle(tb)}))
	metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: fxtest.NewLifecycle(tb)})
}

// EnableMetrics installs the shared test meter provider.
func EnableMetrics(tb testing.TB) {
	tb.Helper()

	metrics.NewMeterProvider(metrics.MeterProviderParams{
		Lifecycle:   fxtest.NewLifecycle(tb),
		Config:      &metrics.Config{},
		Reader:      metrics.NewManualReader(),
		ID:          ID,
		Name:        Name,
		Version:     Version,
		Environment: Environment,
	})
}

// EnableTracer installs the shared test tracer provider.
func EnableTracer(tb testing.TB) {
	tb.Helper()

	require.NoError(tb, tracer.Register(tracer.TracerParams{
		Lifecycle:   fxtest.NewLifecycle(tb),
		Config:      &tracer.Config{Kind: "otlp"},
		ID:          ID,
		Name:        Name,
		Version:     Version,
		Environment: Environment,
	}))
}

// EnableTelemetry installs the shared test meter and tracer providers.
func EnableTelemetry(tb testing.TB) {
	tb.Helper()

	EnableMetrics(tb)
	EnableTracer(tb)
}

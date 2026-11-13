package test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v3/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v3/telemetry/tracer"
	"github.com/alexfalkowski/go-signal/v2/signaltest"
	"github.com/stretchr/testify/require"
)

// ResetTelemetry disables global telemetry state for benchmarks and tests.
func ResetTelemetry(tb testing.TB) {
	tb.Helper()

	require.NoError(tb, tracer.Register(tracer.TracerParams{Lifecycle: signaltest.NewLifecycle(tb)}))
	metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: signaltest.NewLifecycle(tb)})
}

// EnableMetrics installs the shared test meter provider.
//
// This mutates the process-global OpenTelemetry meter provider. Use
// ResetTelemetry or a helper that resets telemetry when a test needs isolation
// from other telemetry assertions in the same process.
func EnableMetrics(tb testing.TB) {
	tb.Helper()

	metrics.NewMeterProvider(metrics.MeterProviderParams{
		Lifecycle:   signaltest.NewLifecycle(tb),
		Config:      &metrics.Config{},
		Reader:      metrics.NewManualReader(),
		ID:          ID,
		Name:        Name,
		Version:     Version,
		Environment: Environment,
	})
}

// EnableTracer installs the shared test tracer provider.
//
// This mutates the process-global OpenTelemetry tracer provider. Use
// ResetTelemetry or EnableIsolatedSpanExporter when a test needs isolated span
// assertions.
func EnableTracer(tb testing.TB) {
	tb.Helper()

	require.NoError(tb, tracer.Register(tracer.TracerParams{
		Lifecycle:   signaltest.NewLifecycle(tb),
		Config:      &tracer.Config{Kind: "otlp", URL: otlpHTTPSURL + "/v1/traces"},
		ID:          ID,
		Name:        Name,
		Version:     Version,
		Environment: Environment,
	}))
}

// EnableTelemetry installs the shared test meter and tracer providers.
//
// This mutates process-global OpenTelemetry meter and tracer providers. Tests
// that assert telemetry state should reset telemetry before and after the
// assertion path.
func EnableTelemetry(tb testing.TB) {
	tb.Helper()

	EnableMetrics(tb)
	EnableTracer(tb)
}

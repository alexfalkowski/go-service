package driver_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

var driverSink driver.Driver

func BenchmarkRedisTelemetry(b *testing.B) {
	bench := func(name string, setup func(testing.TB)) {
		b.Run(name, func(b *testing.B) {
			resetTelemetry(b)
			setup(b)
			defer resetTelemetry(b)

			cfg := &config.Config{
				Kind: "redis",
				Options: map[string]any{
					"url": "redis://localhost:6379",
				},
			}

			lc := noopLifecycle{}

			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				cache, err := driver.NewDriver(driver.DriverParams{
					Lifecycle: lc,
					FS:        test.FS,
					Config:    cfg,
				})
				require.NoError(b, err)
				driverSink = cache
			}
		})
	}

	bench("disabled", func(testing.TB) {})
	bench("metrics", enableMetrics)
	bench("tracer", enableTracer)
	bench("enabled", enableTelemetry)
}

type noopLifecycle struct{}

func (noopLifecycle) Append(di.Hook) {}

func resetTelemetry(tb testing.TB) {
	tb.Helper()

	require.NoError(tb, tracer.Register(tracer.TracerParams{Lifecycle: fxtest.NewLifecycle(tb)}))
	metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: fxtest.NewLifecycle(tb)})
}

func enableMetrics(tb testing.TB) {
	tb.Helper()

	metrics.NewMeterProvider(metrics.MeterProviderParams{
		Lifecycle:   fxtest.NewLifecycle(tb),
		Config:      &metrics.Config{},
		Reader:      metrics.NewManualReader(),
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	})
}

func enableTracer(tb testing.TB) {
	tb.Helper()

	require.NoError(tb, tracer.Register(tracer.TracerParams{
		Lifecycle:   fxtest.NewLifecycle(tb),
		Config:      &tracer.Config{Kind: "otlp"},
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	}))
}

func enableTelemetry(tb testing.TB) {
	tb.Helper()

	enableMetrics(tb)
	enableTracer(tb)
}

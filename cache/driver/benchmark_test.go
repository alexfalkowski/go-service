package driver_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

var driverSink driver.Driver

func BenchmarkRedisTelemetry(b *testing.B) {
	bench := func(name string, setup func(testing.TB)) {
		b.Run(name, func(b *testing.B) {
			test.ResetTelemetry(b)
			setup(b)
			defer test.ResetTelemetry(b)

			cfg := &config.Config{
				Kind: "redis",
				Options: map[string]any{
					"url": "redis://localhost:6379",
				},
			}

			lc := test.NoopLifecycle{}

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
	bench("metrics", test.EnableMetrics)
	bench("tracer", test.EnableTracer)
	bench("enabled", test.EnableTelemetry)
}

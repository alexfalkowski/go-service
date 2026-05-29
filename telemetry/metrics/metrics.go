package metrics

import (
	"github.com/alexfalkowski/go-service/v2/env"
	"go.opentelemetry.io/otel/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// MeterProvider is an alias for metric.MeterProvider.
type MeterProvider = metric.MeterProvider

// Meter is an alias for metric.Meter.
type Meter = metric.Meter

// Registration is an alias for metric.Registration.
type Registration = metric.Registration

// Reader is an alias for sdkmetric.Reader.
type Reader = sdk.Reader

// Metrics is an alias for metricdata.Metrics.
type Metrics = metricdata.Metrics

// ResourceMetrics is an alias for metricdata.ResourceMetrics.
type ResourceMetrics = metricdata.ResourceMetrics

// DataPoint is an alias for metricdata.DataPoint.
type DataPoint[T int64 | float64] = metricdata.DataPoint[T]

// Gauge is an alias for metricdata.Gauge.
type Gauge[T int64 | float64] = metricdata.Gauge[T]

// Sum is an alias for metricdata.Sum.
type Sum[T int64 | float64] = metricdata.Sum[T]

// NewMeter returns a Meter from provider using the service name and version.
//
// The returned meter uses `name` as the instrumentation scope name and `version` as the
// instrumentation scope version (via `metric.WithInstrumentationVersion`).
//
// If provider is nil, NewMeter returns nil.
func NewMeter(name env.Name, version env.Version, provider MeterProvider) Meter {
	if provider == nil {
		return nil
	}

	return provider.Meter(name.String(), metric.WithInstrumentationVersion(version.String()))
}

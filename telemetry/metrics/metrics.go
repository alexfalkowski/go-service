package metrics

import (
	"github.com/alexfalkowski/go-service/v2/env"
	"go.opentelemetry.io/otel/metric"
)

type (

	// MeterProvider is an alias of metric.MeterProvider.
	MeterProvider = metric.MeterProvider

	// Meter is an alias of metric.Meter.
	Meter = metric.Meter
)

// NewMeter for metrics.
func NewMeter(name env.Name, version env.Version, provider MeterProvider) Meter {
	if provider == nil {
		return nil
	}

	return provider.Meter(name.String(), metric.WithInstrumentationVersion(version.String()))
}

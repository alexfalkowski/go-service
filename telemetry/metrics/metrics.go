package metrics

import (
	"github.com/alexfalkowski/go-service/v2/env"
	"go.opentelemetry.io/otel/metric"
)

// MeterProvider is an alias for metric.MeterProvider.
type MeterProvider = metric.MeterProvider

// Meter is an alias for metric.Meter.
type Meter = metric.Meter

// NewMeter returns a Meter from provider using the service name and version.
func NewMeter(name env.Name, version env.Version, provider MeterProvider) Meter {
	if provider == nil {
		return nil
	}

	return provider.Meter(name.String(), metric.WithInstrumentationVersion(version.String()))
}

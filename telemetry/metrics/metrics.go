package metrics

import (
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/runtime"
	"go.opentelemetry.io/otel/metric"
)

// NewMeter for metrics.
func NewMeter(provider metric.MeterProvider, name env.Name) *Meter {
	if provider == nil {
		return nil
	}

	return &Meter{provider.Meter(name.String())}
}

// Meter using otel.
type Meter struct {
	metric.Meter
}

// MustInt64ObservableCounter for metrics.
func (m *Meter) MustInt64ObservableCounter(name, description string) metric.Int64ObservableCounter {
	c, err := m.Int64ObservableCounter(name, metric.WithDescription(description))
	runtime.Must(err)

	return c
}

// MustFloat64ObservableCounter for metrics.
func (m *Meter) MustFloat64ObservableCounter(name, description string) metric.Float64ObservableCounter {
	c, err := m.Float64ObservableCounter(name, metric.WithDescription(description))
	runtime.Must(err)

	return c
}

// MustInt64Counter for metrics.
func (m *Meter) MustInt64Counter(name, description string) metric.Int64Counter {
	c, err := m.Int64Counter(name, metric.WithDescription(description))
	runtime.Must(err)

	return c
}

// MustFloat64Histogram for metrics.
func (m *Meter) MustFloat64Histogram(name, description string) metric.Float64Histogram {
	h, err := m.Float64Histogram(name, metric.WithDescription(description), metric.WithUnit("s"))
	runtime.Must(err)

	return h
}

// MustFloat64Histogram for metrics.
func (m *Meter) MustInt64ObservableGauge(name, description string) metric.Int64ObservableGauge {
	g, err := m.Int64ObservableGauge(name, metric.WithDescription(description))
	runtime.Must(err)

	return g
}

// MustRegisterCallback for metrics.
func (m *Meter) MustRegisterCallback(f metric.Callback, instruments ...metric.Observable) metric.Registration {
	reg, err := m.RegisterCallback(f, instruments...)
	runtime.Must(err)

	return reg
}

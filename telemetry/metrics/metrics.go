package metrics

import (
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type (
	// Float64Histogram is an alias of metric.Float64Histogram.
	Float64Histogram = metric.Float64Histogram

	// Float64ObservableCounter is an alias of metric.Float64ObservableCounter.
	Float64ObservableCounter = metric.Float64ObservableCounter

	// Int64Counter is an alias of metric.Int64Counter.
	Int64Counter = metric.Int64Counter

	// Int64ObservableGauge is an alias of metric.Int64ObservableGauge.
	Int64ObservableGauge = metric.Int64ObservableGauge

	// Int64ObservableCounter is an alias of metric.Int64ObservableCounter.
	Int64ObservableCounter = metric.Int64ObservableCounter

	// Key is an alias of attribute.Key.
	Key = attribute.Key

	// MeasurementOption is an alias of metric.MeasurementOption.
	MeasurementOption = metric.MeasurementOption

	// Observer is an alias of metric.Observer.
	Observer = metric.Observer

	// MeterProvider is an alias of metric.MeterProvider.
	MeterProvider = metric.MeterProvider
)

// WithAttributes is an alias of metric.WithAttributes.
var WithAttributes = metric.WithAttributes

// StringAttr for metrics.
func StringAttr(key, value string) attribute.KeyValue {
	return attribute.Key(key).String(value)
}

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

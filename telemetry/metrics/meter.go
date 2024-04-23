package metrics

import (
	"go.opentelemetry.io/otel/metric"
)

// MustInt64ObservableCounter for metrics.
func MustInt64ObservableCounter(meter metric.Meter, name, description string) metric.Int64ObservableCounter {
	c, err := meter.Int64ObservableCounter(name, metric.WithDescription(description))
	if err != nil {
		panic(err)
	}

	return c
}

// MustFloat64ObservableCounter for metrics.
func MustFloat64ObservableCounter(meter metric.Meter, name, description string) metric.Float64ObservableCounter {
	c, err := meter.Float64ObservableCounter(name, metric.WithDescription(description))
	if err != nil {
		panic(err)
	}

	return c
}

// MustInt64Counter for metrics.
func MustInt64Counter(meter metric.Meter, name, description string) metric.Int64Counter {
	c, err := meter.Int64Counter(name, metric.WithDescription(description))
	if err != nil {
		panic(err)
	}

	return c
}

// MustFloat64Histogram for metrics.
func MustFloat64Histogram(meter metric.Meter, name, description string) metric.Float64Histogram {
	h, err := meter.Float64Histogram(name, metric.WithDescription(description), metric.WithUnit("s"))
	if err != nil {
		panic(err)
	}

	return h
}

// MustFloat64Histogram for metrics.
func MustInt64ObservableGauge(meter metric.Meter, name, description string) metric.Int64ObservableGauge {
	g, err := meter.Int64ObservableGauge(name, metric.WithDescription(description))
	if err != nil {
		panic(err)
	}

	return g
}

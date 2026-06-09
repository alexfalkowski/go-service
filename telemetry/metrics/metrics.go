package metrics

import (
	"github.com/alexfalkowski/go-service/v2/env"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
	"go.opentelemetry.io/otel/metric/noop"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// MeterProvider is an alias for [metric.MeterProvider].
type MeterProvider = metric.MeterProvider

// Meter is an alias for [metric.Meter].
type Meter = metric.Meter

// Registration is an alias for [metric.Registration].
type Registration = metric.Registration

// Callback is an alias for [metric.Callback].
type Callback = metric.Callback

// Observable is an alias for [metric.Observable].
type Observable = metric.Observable

// EmbeddedMeter is an alias for [embedded.Meter].
type EmbeddedMeter = embedded.Meter

// Int64Counter is an alias for [metric.Int64Counter].
type Int64Counter = metric.Int64Counter

// Int64CounterOption is an alias for [metric.Int64CounterOption].
type Int64CounterOption = metric.Int64CounterOption

// Int64UpDownCounter is an alias for [metric.Int64UpDownCounter].
type Int64UpDownCounter = metric.Int64UpDownCounter

// Int64UpDownCounterOption is an alias for [metric.Int64UpDownCounterOption].
type Int64UpDownCounterOption = metric.Int64UpDownCounterOption

// Int64Histogram is an alias for [metric.Int64Histogram].
type Int64Histogram = metric.Int64Histogram

// Int64HistogramOption is an alias for [metric.Int64HistogramOption].
type Int64HistogramOption = metric.Int64HistogramOption

// Int64ObservableCounter is an alias for [metric.Int64ObservableCounter].
type Int64ObservableCounter = metric.Int64ObservableCounter

// Int64ObservableCounterOption is an alias for [metric.Int64ObservableCounterOption].
type Int64ObservableCounterOption = metric.Int64ObservableCounterOption

// Int64ObservableUpDownCounter is an alias for [metric.Int64ObservableUpDownCounter].
type Int64ObservableUpDownCounter = metric.Int64ObservableUpDownCounter

// Int64ObservableUpDownCounterOption is an alias for [metric.Int64ObservableUpDownCounterOption].
type Int64ObservableUpDownCounterOption = metric.Int64ObservableUpDownCounterOption

// Int64ObservableGauge is an alias for [metric.Int64ObservableGauge].
type Int64ObservableGauge = metric.Int64ObservableGauge

// Int64ObservableGaugeOption is an alias for [metric.Int64ObservableGaugeOption].
type Int64ObservableGaugeOption = metric.Int64ObservableGaugeOption

// Int64Gauge is an alias for [metric.Int64Gauge].
type Int64Gauge = metric.Int64Gauge

// Int64GaugeOption is an alias for [metric.Int64GaugeOption].
type Int64GaugeOption = metric.Int64GaugeOption

// Float64Counter is an alias for [metric.Float64Counter].
type Float64Counter = metric.Float64Counter

// Float64CounterOption is an alias for [metric.Float64CounterOption].
type Float64CounterOption = metric.Float64CounterOption

// Float64UpDownCounter is an alias for [metric.Float64UpDownCounter].
type Float64UpDownCounter = metric.Float64UpDownCounter

// Float64UpDownCounterOption is an alias for [metric.Float64UpDownCounterOption].
type Float64UpDownCounterOption = metric.Float64UpDownCounterOption

// Float64Histogram is an alias for [metric.Float64Histogram].
type Float64Histogram = metric.Float64Histogram

// Float64HistogramOption is an alias for [metric.Float64HistogramOption].
type Float64HistogramOption = metric.Float64HistogramOption

// Float64ObservableCounter is an alias for [metric.Float64ObservableCounter].
type Float64ObservableCounter = metric.Float64ObservableCounter

// Float64ObservableCounterOption is an alias for [metric.Float64ObservableCounterOption].
type Float64ObservableCounterOption = metric.Float64ObservableCounterOption

// Float64ObservableUpDownCounter is an alias for [metric.Float64ObservableUpDownCounter].
type Float64ObservableUpDownCounter = metric.Float64ObservableUpDownCounter

// Float64ObservableUpDownCounterOption is an alias for [metric.Float64ObservableUpDownCounterOption].
type Float64ObservableUpDownCounterOption = metric.Float64ObservableUpDownCounterOption

// Float64ObservableGauge is an alias for [metric.Float64ObservableGauge].
type Float64ObservableGauge = metric.Float64ObservableGauge

// Float64ObservableGaugeOption is an alias for [metric.Float64ObservableGaugeOption].
type Float64ObservableGaugeOption = metric.Float64ObservableGaugeOption

// Float64Gauge is an alias for [metric.Float64Gauge].
type Float64Gauge = metric.Float64Gauge

// Float64GaugeOption is an alias for [metric.Float64GaugeOption].
type Float64GaugeOption = metric.Float64GaugeOption

// Reader is an alias for [go.opentelemetry.io/otel/sdk/metric.Reader].
type Reader = sdk.Reader

// Metrics is an alias for [metricdata.Metrics].
type Metrics = metricdata.Metrics

// ResourceMetrics is an alias for [metricdata.ResourceMetrics].
type ResourceMetrics = metricdata.ResourceMetrics

// DataPoint is an alias for [metricdata.DataPoint].
type DataPoint[T int64 | float64] = metricdata.DataPoint[T]

// Gauge is an alias for [metricdata.Gauge].
type Gauge[T int64 | float64] = metricdata.Gauge[T]

// Sum is an alias for [metricdata.Sum].
type Sum[T int64 | float64] = metricdata.Sum[T]

// NewMeter returns a Meter from provider using the service name and version.
//
// The returned meter uses `name` as the instrumentation scope name and `version` as the
// instrumentation scope version (via [metric.WithInstrumentationVersion]).
//
// If provider is nil, NewMeter returns nil.
func NewMeter(name env.Name, version env.Version, provider MeterProvider) Meter {
	if provider == nil {
		return nil
	}

	return provider.Meter(name.String(), metric.WithInstrumentationVersion(version.String()))
}

// NewNoopMeterProvider constructs a no-op meter provider.
func NewNoopMeterProvider() MeterProvider {
	return noop.NewMeterProvider()
}

// GetMeterProvider returns the global OpenTelemetry meter provider.
func GetMeterProvider() MeterProvider {
	return otel.GetMeterProvider()
}

// SetMeterProvider installs the global OpenTelemetry meter provider.
func SetMeterProvider(provider MeterProvider) {
	setMeterProvider(provider, isEnabledProvider(provider))
}

func setMeterProvider(provider MeterProvider, isEnabled bool) {
	otel.SetMeterProvider(provider)
	enabled.Store(isEnabled)
}

func isEnabledProvider(provider MeterProvider) bool {
	if provider == nil {
		return false
	}

	_, ok := provider.(noop.MeterProvider)
	return !ok
}

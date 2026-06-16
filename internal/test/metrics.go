package test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	hm "github.com/alexfalkowski/go-service/v2/transport/http/telemetry/metrics"
	"go.uber.org/fx/fxtest"
)

var errInvalid = errors.New("invalid")

// InvalidMeter is an OpenTelemetry meter test double whose constructors always fail.
type InvalidMeter struct{ metrics.EmbeddedMeter }

// Int64Counter implements [metrics.Meter] and always returns an error.
func (InvalidMeter) Int64Counter(string, ...metrics.Int64CounterOption) (metrics.Int64Counter, error) {
	return nil, errInvalid
}

// Int64UpDownCounter implements [metrics.Meter] and always returns an error.
func (InvalidMeter) Int64UpDownCounter(string, ...metrics.Int64UpDownCounterOption) (metrics.Int64UpDownCounter, error) {
	return nil, errInvalid
}

// Int64Histogram implements [metrics.Meter] and always returns an error.
func (InvalidMeter) Int64Histogram(string, ...metrics.Int64HistogramOption) (metrics.Int64Histogram, error) {
	return nil, errInvalid
}

// Int64ObservableCounter implements [metrics.Meter] and always returns an error.
func (InvalidMeter) Int64ObservableCounter(string, ...metrics.Int64ObservableCounterOption) (metrics.Int64ObservableCounter, error) {
	return nil, errInvalid
}

// Int64ObservableUpDownCounter implements [metrics.Meter] and always returns an error.
func (InvalidMeter) Int64ObservableUpDownCounter(string, ...metrics.Int64ObservableUpDownCounterOption) (metrics.Int64ObservableUpDownCounter, error) {
	return nil, errInvalid
}

// Int64ObservableGauge implements [metrics.Meter] and always returns an error.
func (InvalidMeter) Int64ObservableGauge(string, ...metrics.Int64ObservableGaugeOption) (metrics.Int64ObservableGauge, error) {
	return nil, errInvalid
}

// Int64Gauge implements [metrics.Meter] and always returns an error.
func (InvalidMeter) Int64Gauge(string, ...metrics.Int64GaugeOption) (metrics.Int64Gauge, error) {
	return nil, errInvalid
}

// Float64Counter implements [metrics.Meter] and always returns an error.
func (InvalidMeter) Float64Counter(string, ...metrics.Float64CounterOption) (metrics.Float64Counter, error) {
	return nil, errInvalid
}

// Float64UpDownCounter implements [metrics.Meter] and always returns an error.
func (InvalidMeter) Float64UpDownCounter(string, ...metrics.Float64UpDownCounterOption) (metrics.Float64UpDownCounter, error) {
	return nil, errInvalid
}

// Float64Histogram implements [metrics.Meter] and always returns an error.
func (InvalidMeter) Float64Histogram(string, ...metrics.Float64HistogramOption) (metrics.Float64Histogram, error) {
	return nil, errInvalid
}

// Float64ObservableCounter implements [metrics.Meter] and always returns an error.
func (InvalidMeter) Float64ObservableCounter(string, ...metrics.Float64ObservableCounterOption) (metrics.Float64ObservableCounter, error) {
	return nil, errInvalid
}

// Float64ObservableUpDownCounter implements [metrics.Meter] and always returns an error.
func (InvalidMeter) Float64ObservableUpDownCounter(string, ...metrics.Float64ObservableUpDownCounterOption) (metrics.Float64ObservableUpDownCounter, error) {
	return nil, errInvalid
}

// Float64ObservableGauge implements [metrics.Meter] and always returns an error.
func (InvalidMeter) Float64ObservableGauge(string, ...metrics.Float64ObservableGaugeOption) (metrics.Float64ObservableGauge, error) {
	return nil, errInvalid
}

// Float64Gauge implements [metrics.Meter] and always returns an error.
func (InvalidMeter) Float64Gauge(string, ...metrics.Float64GaugeOption) (metrics.Float64Gauge, error) {
	return nil, errInvalid
}

// RegisterCallback implements [metrics.Meter] and always returns an error.
func (InvalidMeter) RegisterCallback(metrics.Callback, ...metrics.Observable) (metrics.Registration, error) {
	return nil, errInvalid
}

// NewOTLPMeter returns a meter backed by the shared OTLP metrics config.
func NewOTLPMeter(lc di.Lifecycle) (metrics.Meter, error) {
	return NewMeter(lc, NewOTLPMetricsConfig())
}

// NewPrometheusMeter returns a meter backed by the shared Prometheus metrics config.
func NewPrometheusMeter(lc di.Lifecycle) (metrics.Meter, error) {
	return NewMeter(lc, NewPrometheusMetricsConfig())
}

// NewMeter returns a repository meter scoped to the shared test name and version.
func NewMeter(lc di.Lifecycle, c *metrics.Config) (metrics.Meter, error) {
	provider, err := NewMeterProvider(lc, c)
	if err != nil {
		return nil, err
	}

	return metrics.NewMeter(Name, Version, provider), nil
}

// NewOTLPMeterProvider returns a meter provider backed by the shared OTLP metrics config.
func NewOTLPMeterProvider(lc di.Lifecycle) (metrics.MeterProvider, error) {
	return NewMeterProvider(lc, NewOTLPMetricsConfig())
}

// NewPrometheusMeterProvider returns a meter provider backed by the shared Prometheus metrics config.
func NewPrometheusMeterProvider(lc di.Lifecycle) (metrics.MeterProvider, error) {
	return NewMeterProvider(lc, NewPrometheusMetricsConfig())
}

// NewMeterProvider creates a meter provider with a reader registered on the supplied lifecycle.
func NewMeterProvider(lc di.Lifecycle, config *metrics.Config) (metrics.MeterProvider, error) {
	r, err := metrics.NewReader(lc, Name, config)
	if err != nil {
		return nil, err
	}

	params := metrics.MeterProviderParams{
		Lifecycle:   lc,
		Config:      config,
		Reader:      r,
		Environment: Environment,
		Version:     Version,
		Name:        Name,
	}

	return metrics.NewMeterProvider(params), nil
}

// EnableMetricsReader installs the shared test meter provider and returns its manual reader.
//
// It resets process-global telemetry before installation and again with
// tb.Cleanup, so the returned reader can be used for isolated metric
// assertions.
func EnableMetricsReader(tb testing.TB) metrics.Reader {
	tb.Helper()

	ResetTelemetry(tb)
	tb.Cleanup(func() {
		ResetTelemetry(tb)
	})

	reader := metrics.NewManualReader()
	metrics.NewMeterProvider(metrics.MeterProviderParams{
		Lifecycle:   fxtest.NewLifecycle(tb),
		Config:      &metrics.Config{},
		Reader:      reader,
		ID:          ID,
		Name:        Name,
		Version:     Version,
		Environment: Environment,
	})

	return reader
}

func meter(lc di.Lifecycle, mux *http.ServeMux, os *worldOpts) (metrics.Meter, error) {
	if os.telemetry == "otlp" {
		return NewMeter(lc, NewOTLPMetricsConfig())
	}

	config := NewPrometheusMetricsConfig()
	hm.Register(Name, config, mux)

	return NewMeter(lc, config)
}

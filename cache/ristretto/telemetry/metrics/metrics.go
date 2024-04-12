package metrics

import (
	"context"

	"github.com/alexfalkowski/go-service/cache/ristretto"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
	m "github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/version"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
)

// Meter for metrics.
type Meter metric.Meter

// Params for metrics.
type Params struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Config      *m.Config
	Environment env.Environment
	Version     version.Version
}

// NewMeter for metrics.
func NewMeter(params Params) (Meter, error) {
	return m.NewMeter(params.Lifecycle, "http", params.Environment, params.Version, params.Config)
}

// Register for metrics.
func Register(cache ristretto.Cache, version version.Version, meter Meter) error {
	opts := metric.WithAttributes(
		attribute.Key("name").String(os.ExecutableName()),
		attribute.Key("version").String(string(version)),
	)

	hits, err := meter.Int64ObservableCounter("ristretto_hits_total", metric.WithDescription("The number of hits in the cache."))
	if err != nil {
		return err
	}

	misses, err := meter.Int64ObservableCounter("ristretto_misses_total", metric.WithDescription("The number of misses in the cache."))
	if err != nil {
		return err
	}

	m := &metrics{cache: cache, opts: opts, hit: hits, miss: misses}

	meter.RegisterCallback(m.callback, hits, misses)

	return nil
}

type metrics struct {
	cache ristretto.Cache
	opts  metric.MeasurementOption

	hit  metric.Int64ObservableCounter
	miss metric.Int64ObservableCounter
}

func (m *metrics) callback(_ context.Context, o metric.Observer) error {
	o.ObserveInt64(m.hit, int64(m.cache.Hits()), m.opts)
	o.ObserveInt64(m.miss, int64(m.cache.Misses()), m.opts)

	return nil
}

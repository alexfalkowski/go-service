package metrics

import (
	"context"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	"github.com/dgraph-io/ristretto"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Register for metrics.
func Register(cache *ristretto.Cache, version version.Version, meter metric.Meter) error {
	opts := metric.WithAttributes(
		attribute.Key("name").String(os.ExecutableName()),
		attribute.Key("version").String(string(version)),
	)

	hits, err := meter.Float64ObservableCounter("ristretto_hits_total", metric.WithDescription("The number of hits in the cache."))
	if err != nil {
		return err
	}

	misses, err := meter.Float64ObservableCounter("ristretto_misses_total", metric.WithDescription("The number of misses in the cache."))
	if err != nil {
		return err
	}

	m := &metrics{cache: cache, opts: opts, hit: hits, miss: misses}

	meter.RegisterCallback(m.callback, hits, misses)

	return nil
}

type metrics struct {
	cache *ristretto.Cache
	opts  metric.MeasurementOption

	hit  metric.Float64ObservableCounter
	miss metric.Float64ObservableCounter
}

func (m *metrics) callback(_ context.Context, o metric.Observer) error {
	stats := m.cache.Metrics

	o.ObserveFloat64(m.hit, float64(stats.Hits()), m.opts)
	o.ObserveFloat64(m.miss, float64(stats.Misses()), m.opts)

	return nil
}

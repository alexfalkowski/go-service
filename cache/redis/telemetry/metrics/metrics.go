package metrics

import (
	"context"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	"github.com/go-redis/cache/v8"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Register for metrics.
func Register(cache *cache.Cache, version version.Version, meter metric.Meter) error {
	opts := metric.WithAttributes(
		attribute.Key("name").String(os.ExecutableName()),
		attribute.Key("version").String(string(version)),
	)

	hits, err := meter.Int64ObservableCounter("redis_hits_total", metric.WithDescription("The number of hits in the cache."))
	if err != nil {
		return err
	}

	misses, err := meter.Int64ObservableCounter("redis_misses_total", metric.WithDescription("The number of misses in the cache."))
	if err != nil {
		return err
	}

	m := &metrics{cache: cache, opts: opts, hit: hits, miss: misses}

	meter.RegisterCallback(m.callback, hits, misses)

	return nil
}

type metrics struct {
	cache *cache.Cache
	opts  metric.MeasurementOption

	hit  metric.Int64ObservableCounter
	miss metric.Int64ObservableCounter
}

func (m *metrics) callback(_ context.Context, o metric.Observer) error {
	stats := m.cache.Stats()

	o.ObserveInt64(m.hit, int64(stats.Hits), m.opts)
	o.ObserveInt64(m.miss, int64(stats.Misses), m.opts)

	return nil
}

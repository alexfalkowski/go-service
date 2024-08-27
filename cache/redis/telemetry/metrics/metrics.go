package metrics

import (
	"context"

	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/go-redis/cache/v9"
	"go.opentelemetry.io/otel/metric"
)

// Register for metrics.
func Register(cache *cache.Cache, meter metric.Meter) {
	hits := metrics.MustInt64ObservableCounter(meter, "redis_hits_total", "The number of hits in the cache.")
	misses := metrics.MustInt64ObservableCounter(meter, "redis_misses_total", "The number of misses in the cache.")
	m := &ms{cache: cache, hit: hits, miss: misses}

	meter.RegisterCallback(m.callback, hits, misses)
}

type ms struct {
	cache *cache.Cache
	hit   metric.Int64ObservableCounter
	miss  metric.Int64ObservableCounter
}

//nolint:gosec
func (m *ms) callback(_ context.Context, o metric.Observer) error {
	stats := m.cache.Stats()
	if stats != nil {
		o.ObserveInt64(m.hit, int64(stats.Hits))
		o.ObserveInt64(m.miss, int64(stats.Misses))
	}

	return nil
}

package metrics

import (
	"context"

	"github.com/alexfalkowski/go-service/ristretto"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"go.opentelemetry.io/otel/metric"
)

// Register for metrics.
func Register(cache ristretto.Cache, meter metric.Meter) {
	hits := metrics.MustInt64ObservableCounter(meter, "ristretto_hits_total", "The number of hits in the cache.")
	misses := metrics.MustInt64ObservableCounter(meter, "ristretto_misses_total", "The number of misses in the cache.")
	m := &ms{cache: cache, hit: hits, miss: misses}

	meter.RegisterCallback(m.callback, hits, misses)
}

type ms struct {
	cache ristretto.Cache
	hit   metric.Int64ObservableCounter
	miss  metric.Int64ObservableCounter
}

//nolint:gosec
func (m *ms) callback(_ context.Context, o metric.Observer) error {
	o.ObserveInt64(m.hit, int64(m.cache.Hits()))
	o.ObserveInt64(m.miss, int64(m.cache.Misses()))

	return nil
}

package metrics

import (
	"context"

	"github.com/alexfalkowski/go-service/cache/ristretto"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/version"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Register for metrics.
func Register(cache ristretto.Cache, version version.Version, meter metric.Meter) {
	opts := metric.WithAttributes(
		attribute.Key("name").String(os.ExecutableName()),
		attribute.Key("version").String(string(version)),
	)

	hits := metrics.MustInt64ObservableCounter(meter, "ristretto_hits_total", "The number of hits in the cache.")
	misses := metrics.MustInt64ObservableCounter(meter, "ristretto_misses_total", "The number of misses in the cache.")
	m := &ms{cache: cache, opts: opts, hit: hits, miss: misses}

	meter.RegisterCallback(m.callback, hits, misses)
}

type ms struct {
	cache ristretto.Cache
	opts  metric.MeasurementOption

	hit  metric.Int64ObservableCounter
	miss metric.Int64ObservableCounter
}

func (m *ms) callback(_ context.Context, o metric.Observer) error {
	o.ObserveInt64(m.hit, int64(m.cache.Hits()), m.opts)
	o.ObserveInt64(m.miss, int64(m.cache.Misses()), m.opts)

	return nil
}

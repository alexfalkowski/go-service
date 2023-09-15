package telemetry

import (
	"context"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	"github.com/go-redis/cache/v8"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

// Register for telemetry.
func RegisterMetrics(lc fx.Lifecycle, metrics *Metrics) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return prometheus.Register(metrics)
		},
		OnStop: func(ctx context.Context) error {
			prometheus.Unregister(metrics)

			return nil
		},
	})
}

// Metrics implements the prometheus.Collector interface.
type Metrics struct {
	cache  *cache.Cache
	hits   *prometheus.Desc
	misses *prometheus.Desc
}

// NewMetrics for telemetry.
func NewMetrics(cache *cache.Cache, version version.Version) *Metrics {
	labels := prometheus.Labels{"name": os.ExecutableName(), "version": string(version)}

	return &Metrics{
		cache: cache,
		hits: prometheus.NewDesc(
			"redis_hits_total",
			"The number of hits in the cache.",
			nil,
			labels,
		),
		misses: prometheus.NewDesc(
			"redis_misses_total",
			"The number of misses in the cache.",
			nil,
			labels,
		),
	}
}

// Describe implements the prometheus.Collector interface.
func (c *Metrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.hits
	ch <- c.misses
}

// Collect implements the prometheus.Collector interface.
func (c *Metrics) Collect(ch chan<- prometheus.Metric) {
	stats := c.cache.Stats()

	ch <- prometheus.MustNewConstMetric(
		c.hits,
		prometheus.CounterValue,
		float64(stats.Hits),
	)
	ch <- prometheus.MustNewConstMetric(
		c.misses,
		prometheus.CounterValue,
		float64(stats.Misses),
	)
}

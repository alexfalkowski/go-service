package telemetry

import (
	"context"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	"github.com/dgraph-io/ristretto"
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
	cache  *ristretto.Cache
	hits   *prometheus.Desc
	misses *prometheus.Desc
}

// NewMetrics for prometheus.
func NewMetrics(cache *ristretto.Cache, version version.Version) *Metrics {
	labels := prometheus.Labels{"name": os.ExecutableName(), "version": string(version)}

	return &Metrics{
		cache: cache,
		hits: prometheus.NewDesc(
			"ristretto_hits_total",
			"The number of hits in the cache.",
			nil,
			labels,
		),
		misses: prometheus.NewDesc(
			"ristretto_misses_total",
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
	metrics := c.cache.Metrics

	ch <- prometheus.MustNewConstMetric(
		c.hits,
		prometheus.CounterValue,
		float64(metrics.Hits()),
	)
	ch <- prometheus.MustNewConstMetric(
		c.misses,
		prometheus.CounterValue,
		float64(metrics.Misses()),
	)
}

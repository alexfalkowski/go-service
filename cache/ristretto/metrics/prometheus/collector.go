package prometheus

import (
	"github.com/dgraph-io/ristretto"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "go_ristretto_stats"
	subsystem = "cache"
)

// StatsCollector implements the prometheus.Collector interface.
type StatsCollector struct {
	cache *ristretto.Cache

	// descriptions of exported metrics
	hitsDesc   *prometheus.Desc
	missesDesc *prometheus.Desc
}

// NewStatsCollector creates a new StatsCollector.
func NewStatsCollector(cacheName string, cache *ristretto.Cache) *StatsCollector {
	labels := prometheus.Labels{"cache_name": cacheName}

	return &StatsCollector{
		cache: cache,
		hitsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "hits"),
			"The number of hits in the cache.",
			nil,
			labels,
		),
		missesDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "misses"),
			"The number of misses in the cache.",
			nil,
			labels,
		),
	}
}

// Describe implements the prometheus.Collector interface.
func (c StatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.hitsDesc
	ch <- c.missesDesc
}

// Collect implements the prometheus.Collector interface.
func (c StatsCollector) Collect(ch chan<- prometheus.Metric) {
	metrics := c.cache.Metrics

	ch <- prometheus.MustNewConstMetric(
		c.hitsDesc,
		prometheus.CounterValue,
		float64(metrics.Hits()),
	)
	ch <- prometheus.MustNewConstMetric(
		c.missesDesc,
		prometheus.CounterValue,
		float64(metrics.Misses()),
	)
}

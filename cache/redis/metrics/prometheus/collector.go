package prometheus

import (
	"github.com/go-redis/cache/v8"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "go_redis_stats"
	subsystem = "cache"
)

// StatsGetter is an interface that gets *cache.Stats.
type StatsGetter interface {
	Stats() *cache.Stats
}

// StatsCollector implements the prometheus.Collector interface.
type StatsCollector struct {
	sg StatsGetter

	// descriptions of exported metrics
	hitsDesc   *prometheus.Desc
	missesDesc *prometheus.Desc
}

// NewStatsCollector creates a new StatsCollector.
func NewStatsCollector(cacheName string, sg StatsGetter) *StatsCollector {
	labels := prometheus.Labels{"cache_name": cacheName}

	return &StatsCollector{
		sg: sg,
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
	stats := c.sg.Stats()

	ch <- prometheus.MustNewConstMetric(
		c.hitsDesc,
		prometheus.CounterValue,
		float64(stats.Hits),
	)
	ch <- prometheus.MustNewConstMetric(
		c.missesDesc,
		prometheus.CounterValue,
		float64(stats.Misses),
	)
}

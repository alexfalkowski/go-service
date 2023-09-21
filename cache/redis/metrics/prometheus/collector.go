package prometheus

import (
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	"github.com/go-redis/cache/v8"
	"github.com/prometheus/client_golang/prometheus"
)

// StatsCollector implements the prometheus.Collector interface.
type StatsCollector struct {
	cache *cache.Cache

	// descriptions of exported metrics
	hitsDesc   *prometheus.Desc
	missesDesc *prometheus.Desc
}

// NewStatsCollector for prometheus.
func NewStatsCollector(cache *cache.Cache, version version.Version) *StatsCollector {
	labels := prometheus.Labels{"name": os.ExecutableName(), "version": string(version)}

	return &StatsCollector{
		cache: cache,
		hitsDesc: prometheus.NewDesc(
			"redis_hits_total",
			"The number of hits in the cache.",
			nil,
			labels,
		),
		missesDesc: prometheus.NewDesc(
			"redis_misses_total",
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
	stats := c.cache.Stats()

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

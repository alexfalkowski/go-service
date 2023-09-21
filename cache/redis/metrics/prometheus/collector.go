package prometheus

import (
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	"github.com/go-redis/cache/v8"
	"github.com/prometheus/client_golang/prometheus"
)

// Collector implements the prometheus.Collector interface.
type Collector struct {
	cache  *cache.Cache
	hits   *prometheus.Desc
	misses *prometheus.Desc
}

// NewCollector for prometheus.
func NewCollector(cache *cache.Cache, version version.Version) *Collector {
	labels := prometheus.Labels{"name": os.ExecutableName(), "version": string(version)}

	return &Collector{
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
func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.hits
	ch <- c.misses
}

// Collect implements the prometheus.Collector interface.
func (c Collector) Collect(ch chan<- prometheus.Metric) {
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

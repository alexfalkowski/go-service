package prometheus

import (
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	"github.com/dgraph-io/ristretto"
	"github.com/prometheus/client_golang/prometheus"
)

// Collector implements the prometheus.Collector interface.
type Collector struct {
	cache  *ristretto.Cache
	hits   *prometheus.Desc
	misses *prometheus.Desc
}

// NewCollector for prometheus.
func NewCollector(cache *ristretto.Cache, version version.Version) *Collector {
	labels := prometheus.Labels{"name": os.ExecutableName(), "version": string(version)}

	return &Collector{
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
func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.hits
	ch <- c.misses
}

// Collect implements the prometheus.Collector interface.
func (c Collector) Collect(ch chan<- prometheus.Metric) {
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

package prometheus

import (
	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/go-redis/cache/v8"
	"github.com/prometheus/client_golang/prometheus"
)

// Register for prometheus.
func Register(cfg *config.Config, cache *cache.Cache) error {
	collector := NewStatsCollector(cfg.AppName, cache)

	return prometheus.Register(collector)
}

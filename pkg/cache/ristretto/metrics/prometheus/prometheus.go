package prometheus

import (
	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/dgraph-io/ristretto"
	"github.com/prometheus/client_golang/prometheus"
)

// Register for prometheus.
func Register(cfg *config.Config, cache *ristretto.Cache) error {
	collector := NewStatsCollector(cfg.AppName, cache)

	return prometheus.Register(collector)
}

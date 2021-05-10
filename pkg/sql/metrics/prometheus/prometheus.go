package prometheus

import (
	"database/sql"

	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/dlmiddlecote/sqlstats"
	"github.com/prometheus/client_golang/prometheus"
)

// Register for prometheus.
func Register(cfg *config.Config, db *sql.DB) error {
	collector := sqlstats.NewStatsCollector(cfg.AppName, db)

	return prometheus.Register(collector)
}

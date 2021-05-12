package prometheus

import (
	"context"
	"database/sql"

	"github.com/dlmiddlecote/sqlstats"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

// Register for prometheus.
func Register(lc fx.Lifecycle, appName string, db *sql.DB) {
	collector := sqlstats.NewStatsCollector(appName, db)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return prometheus.Register(collector)
		},
		OnStop: func(ctx context.Context) error {
			prometheus.Unregister(collector)

			return nil
		},
	})
}

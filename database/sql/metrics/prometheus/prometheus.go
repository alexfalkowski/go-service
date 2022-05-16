package prometheus

import (
	"context"
	"database/sql"

	"github.com/alexfalkowski/go-service/version"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

// Register for prometheus.
func Register(lc fx.Lifecycle, name string, db *sql.DB, version version.Version) {
	collector := NewStatsCollector(name, db, version)

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

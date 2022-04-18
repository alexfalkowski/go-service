package pg

import (
	"context"
	"database/sql"

	"github.com/alexfalkowski/go-service/database/sql/metrics/prometheus"
	"github.com/alexfalkowski/go-service/os"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/fx"
)

// NewDB for SQL.
func NewDB(lc fx.Lifecycle, cfg *Config) (*sql.DB, error) {
	config, err := pgx.ParseConfig(cfg.URL)
	if err != nil {
		return nil, err
	}

	db := stdlib.OpenDB(*config)

	prometheus.Register(lc, os.ExecutableName(), db)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db, nil
}

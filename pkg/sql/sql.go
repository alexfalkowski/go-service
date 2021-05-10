package sql

import (
	"context"
	"database/sql"

	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/alexfalkowski/go-service/pkg/sql/metrics/prometheus"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/fx"
)

// NewDB for SQL.
func NewDB(lc fx.Lifecycle, cfg *config.Config) (*sql.DB, error) {
	config, err := pgx.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	db := stdlib.OpenDB(*config)

	if err := prometheus.Register(cfg, db); err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db, nil
}

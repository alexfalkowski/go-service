package sql

import (
	"context"
	"database/sql"

	"github.com/alexfalkowski/go-service/pkg/config"
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

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db, nil
}

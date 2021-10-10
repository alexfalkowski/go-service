package pg

import (
	"context"
	"database/sql"

	"github.com/alexfalkowski/go-service/pkg/os"
	"github.com/alexfalkowski/go-service/pkg/sql/metrics/prometheus"
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

	name, err := os.ExecutableName()
	if err != nil {
		return nil, err
	}

	prometheus.Register(lc, name, db)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db, nil
}

package pg

import (
	"context"
	"database/sql"

	"github.com/alexfalkowski/go-service/database/sql/metrics/prometheus"
	"github.com/alexfalkowski/go-service/version"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/fx"
)

// DBParams for SQL.
type DBParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *Config
	Version   version.Version
}

// NewDB for SQL.
func NewDB(params DBParams) (*sql.DB, error) {
	config, err := pgx.ParseConfig(params.Config.URL)
	if err != nil {
		return nil, err
	}

	db := stdlib.OpenDB(*config)

	prometheus.Register(params.Lifecycle, db, params.Version)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db, nil
}

package pg

import (
	"sync"

	"github.com/alexfalkowski/go-service/database/sql/driver"
	"github.com/alexfalkowski/go-service/database/sql/pg/otel"
	"github.com/alexfalkowski/go-service/version"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/linxGnu/mssqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// DBParams for PostgreSQL.
type DBParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *Config
	Version   version.Version
}

// Open for PostgreSQL.
func Open(params DBParams) (*mssqlx.DBs, error) {
	return driver.Open(params.Lifecycle, "pg", params.Config.Config, params.Version)
}

var once sync.Once

// Register the driver for PostgreSQL.
func Register(tracer otel.Tracer, logger *zap.Logger) {
	once.Do(func() {
		driver.Register("pg", stdlib.GetDefaultDriver(), tracer, logger)
	})
}

package pg

import (
	"sync"

	"github.com/alexfalkowski/go-service/database/sql/driver"
	"github.com/alexfalkowski/go-service/database/sql/pg/telemetry/tracer"
	"github.com/alexfalkowski/go-service/version"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/linxGnu/mssqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// OpenParams for pg.
type OpenParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *Config
	Version   version.Version
}

// Open for pg.
func Open(params OpenParams) (*mssqlx.DBs, error) {
	return driver.Open(params.Lifecycle, "pg", params.Config.Config, params.Version)
}

var once sync.Once

// Register the driver for pg.
func Register(tracer tracer.Tracer, logger *zap.Logger) {
	once.Do(func() {
		driver.Register("pg", stdlib.GetDefaultDriver(), tracer, logger)
	})
}

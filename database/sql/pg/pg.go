package pg

import (
	"sync"

	"github.com/alexfalkowski/go-service/database/sql/driver"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/linxGnu/mssqlx"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// OpenParams for pg.
type OpenParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *Config
}

// Open for pg.
func Open(params OpenParams) (*mssqlx.DBs, error) {
	if !IsEnabled(params.Config) {
		return &mssqlx.DBs{}, nil
	}

	return driver.Open(params.Lifecycle, "pg", params.Config.Config)
}

var once sync.Once

// Register the driver for pg.
func Register(tracer trace.Tracer, logger *zap.Logger) {
	once.Do(func() {
		driver.Register("pg", stdlib.GetDefaultDriver(), tracer, logger)
	})
}

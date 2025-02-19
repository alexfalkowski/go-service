package pg

import (
	"sync"

	"github.com/alexfalkowski/go-service/database/sql/driver"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/linxGnu/mssqlx"
	"go.uber.org/fx"
)

// Open for pg.
func Open(lc fx.Lifecycle, cfg *Config) (*mssqlx.DBs, error) {
	if !IsEnabled(cfg) {
		return &mssqlx.DBs{}, nil
	}

	return driver.Open(lc, "pg", cfg.Config)
}

var once sync.Once

// Register the driver for pg.
func Register(tracer *tracer.Tracer, logger *logger.Logger) {
	once.Do(func() {
		driver.Register("pg", stdlib.GetDefaultDriver(), tracer, logger)
	})
}

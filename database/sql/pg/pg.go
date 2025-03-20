package pg

import (
	"github.com/alexfalkowski/go-service/database/sql/driver"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	pgx "github.com/jackc/pgx/v5/stdlib"
	"github.com/linxGnu/mssqlx"
	"go.uber.org/fx"
)

// Register for pg.
func Register(trace *tracer.Tracer, log *logger.Logger) {
	_ = driver.Register("pg", driver.New("pg", pgx.GetDefaultDriver(), trace, log))
}

// Open for pg.
func Open(lc fx.Lifecycle, cfg *Config) (*mssqlx.DBs, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	return driver.Open(lc, "pg", cfg.Config)
}

package pg

import (
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	pgx "github.com/jackc/pgx/v5/stdlib"
	"github.com/linxGnu/mssqlx"
	"go.uber.org/fx"
)

// Register for pg.
func Register(trace *tracer.Tracer, log *logger.Logger) {
	_ = driver.Register("pg", driver.New("pg", pgx.GetDefaultDriver(), trace, log))
}

// Open for pg.
func Open(lc fx.Lifecycle, fs *os.FS, cfg *Config) (*mssqlx.DBs, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	return driver.Open(lc, "pg", fs, cfg.Config)
}

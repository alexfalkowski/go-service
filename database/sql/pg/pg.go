package pg

import (
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	pgx "github.com/jackc/pgx/v5/stdlib"
	"github.com/linxGnu/mssqlx"
)

// Register for pg.
func Register(trace *tracer.Tracer, log *logger.Logger) {
	_ = driver.Register("pg", pgx.GetDefaultDriver())
}

// Open for pg.
func Open(lc di.Lifecycle, fs *os.FS, cfg *Config) (*mssqlx.DBs, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	return driver.Open(lc, "pg", fs, cfg.Config)
}

package pg

import (
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/os"
	pgx "github.com/jackc/pgx/v5/stdlib"
	"github.com/linxGnu/mssqlx"
)

// Register registers the pgx `database/sql` driver under the name "pg" and enables
// OpenTelemetry driver instrumentation.
func Register() {
	_ = driver.Register("pg", pgx.GetDefaultDriver())
}

// Open opens PostgreSQL master/slave connection pools.
//
// If cfg is disabled, it returns (nil, nil).
func Open(lc di.Lifecycle, fs *os.FS, cfg *Config) (*mssqlx.DBs, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	return driver.Open(lc, "pg", fs, cfg.Config)
}

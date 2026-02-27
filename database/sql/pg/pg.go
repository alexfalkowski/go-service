package pg

import (
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/os"
	pgx "github.com/jackc/pgx/v5/stdlib"
	"github.com/linxGnu/mssqlx"
)

// Register registers the pgx stdlib `database/sql` driver under the name "pg".
//
// The registration is performed via `database/sql/driver.Register`, which wraps the underlying
// driver with OpenTelemetry instrumentation (via otelsql.WrapDriver). The returned error from
// registration is intentionally ignored.
//
// Register is typically called during process initialization via DI wiring (see `pg.Module`).
func Register() {
	_ = driver.Register("pg", pgx.GetDefaultDriver())
}

// Open opens PostgreSQL master/slave connection pools.
//
// Disabled behavior: if cfg is nil/disabled, Open returns (nil, nil).
//
// Enabled behavior: Open delegates to the shared SQL driver helper to:
//   - resolve master and replica DSNs (expressed as go-service "source strings"),
//   - connect using the previously registered driver name "pg",
//   - register OpenTelemetry DB stats metrics, and
//   - apply connection pool limits/lifetime.
//
// The returned pools are closed on Fx lifecycle stop (via an OnStop hook).
func Open(lc di.Lifecycle, fs *os.FS, cfg *Config) (*mssqlx.DBs, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	return driver.Open(lc, "pg", fs, cfg.Config)
}

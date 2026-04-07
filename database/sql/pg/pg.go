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

// Connect opens PostgreSQL master/slave connection pools.
//
// Disabled behavior: if cfg is nil/disabled, Connect returns (nil, nil).
//
// Enabled behavior: Connect delegates to the shared SQL driver helper to:
//   - resolve master and replica DSNs (expressed as go-service "source strings"),
//   - connect using the previously registered driver name "pg",
//   - register OpenTelemetry DB stats metrics, and
//   - apply connection pool limits/lifetime.
func Connect(fs *os.FS, cfg *Config) (*mssqlx.DBs, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	return driver.Connect("pg", fs, cfg.Config)
}

// Open opens PostgreSQL master/slave connection pools.
//
// Open preserves PostgreSQL's nil/disabled config semantics and then delegates
// connection lifecycle ownership to the shared SQL driver helper.
func Open(lc di.Lifecycle, fs *os.FS, cfg *Config) (*mssqlx.DBs, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	return driver.Open(lc, "pg", fs, cfg.Config)
}

package pg

import (
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/database/sql/telemetry"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	pgx "github.com/jackc/pgx/v5/stdlib"
)

// Register registers the pgx stdlib `database/sql` driver under the name "pg".
//
// The registration is performed via `database/sql/driver.Register`, which wraps
// the underlying driver with OpenTelemetry instrumentation via
// `database/sql/telemetry`. The returned error from registration is
// intentionally ignored.
//
// Register is typically called during process initialization via DI wiring (see `pg.Module`).
func Register() {
	_ = driver.Register("pg", pgx.GetDefaultDriver(), options()...)
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
//
// The returned type wraps the upstream master/slave pool collection and is
// aliased by the root `database/sql` package as `sql.DBs` for higher-level
// callers.
func Connect(fs *os.FS, cfg *Config) (*driver.DBs, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	return driver.Connect("pg", fs, cfg.Config, options()...)
}

// Open opens PostgreSQL master/slave connection pools.
//
// Open preserves PostgreSQL's nil/disabled config semantics and then delegates
// connection lifecycle ownership to the shared SQL driver helper.
//
// The returned type is the same go-service DBs wrapper returned by Connect.
func Open(lc di.Lifecycle, fs *os.FS, cfg *Config) (*driver.DBs, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	return driver.Open(lc, "pg", fs, cfg.Config, options()...)
}

func options() []telemetry.Option {
	return []telemetry.Option{telemetry.WithAttributes(attributes.DBSystemNamePostgreSQL)}
}

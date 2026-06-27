package pg

import (
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/database/sql/telemetry"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	pgx "github.com/jackc/pgx/v5/stdlib"
)

// Register registers the pgx stdlib [database/sql] driver under the name "pg".
//
// The registration is performed via [github.com/alexfalkowski/go-service/v2/database/sql/driver.Register], which wraps
// the underlying driver with OpenTelemetry instrumentation via
// [github.com/alexfalkowski/go-service/v2/database/sql/telemetry] when tracing or metrics are enabled.
// The returned error from registration is intentionally ignored.
//
// Register is typically called during process initialization via DI wiring (see [Module]).
func Register() {
	_ = driver.Register("pg", pgx.GetDefaultDriver(), options()...)
}

// Connect opens PostgreSQL writer/reader connection pools.
//
// Disabled behavior: if cfg is nil/disabled, Connect returns (nil, nil).
//
// Enabled behavior: Connect delegates to the shared SQL driver helper to:
//   - resolve writer and replica DSNs (expressed as go-service "source strings"),
//   - require at least one non-empty writer or replica DSN,
//   - create pool handles using the previously registered driver name "pg",
//   - register OpenTelemetry DB stats metrics when metrics are enabled, and
//   - apply connection pool limits/lifetime.
//
// Pool creation follows the database/sql Open contract and does not ping the
// database or verify network reachability. Call the returned DBs Ping helpers or
// register a DB health checker when readiness must prove connectivity.
//
// The returned type is aliased by the root [github.com/alexfalkowski/go-service/v2/database/sql] package as
// [github.com/alexfalkowski/go-service/v2/database/sql.DBs] for higher-level callers.
func Connect(fs *os.FS, cfg *Config) (*driver.DBs, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	return driver.Connect("pg", fs, cfg.Config, options()...)
}

// Open opens PostgreSQL writer/reader connection pools.
//
// Open preserves PostgreSQL's nil/disabled config semantics and then delegates
// connection lifecycle ownership to the shared SQL driver helper.
//
// The returned type is the same go-service [driver.DBs] wrapper returned by [Connect].
func Open(lc di.Lifecycle, fs *os.FS, cfg *Config) (*driver.DBs, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	return driver.Open(lc, "pg", fs, cfg.Config, options()...)
}

func options() []telemetry.Option {
	return []telemetry.Option{telemetry.WithAttributes(attributes.DBSystemNamePostgreSQL)}
}

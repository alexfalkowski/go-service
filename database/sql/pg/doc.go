// Package pg provides PostgreSQL ([github.com/alexfalkowski/go-service/v2/database/sql]) wiring and helpers for go-service.
//
// This package integrates the [github.com/jackc/pgx/v5/stdlib] PostgreSQL driver with go-service SQL wiring by:
//
//   - registering the pgx stdlib [database/sql] driver under the name "pg", with OpenTelemetry instrumentation
//     when tracing or metrics are enabled, and
//   - providing an [Open] constructor that creates writer/reader connection pools using the shared SQL driver helpers.
//
// # Configuration and enablement
//
// PostgreSQL configuration is optional. By convention, a nil *[Config] (or nil embedded config) is treated as
// "disabled", and constructors such as [Open] return (nil, nil) when disabled.
// Enabled PostgreSQL configuration must provide at least one non-empty writer or replica DSN.
// PostgreSQL DSN options, including TLS and sslmode settings, are supplied in the configured DSN and
// passed through to pgx without repository-level rewriting or security policy enforcement.
//
// # Writer/reader pools
//
// [Open] resolves writer and replica DSNs from configuration (DSNs are
// expressed as go-service "source strings"), connects using the shared
// writer/reader pool abstraction used by the repository, applies pool settings
// (max lifetime/open/idle), and registers OpenTelemetry DB stats metrics when
// metrics are enabled.
//
// The package returns [github.com/alexfalkowski/go-service/v2/database/sql/driver.DBs], which is aliased by the root
// [github.com/alexfalkowski/go-service/v2/database/sql] package as [github.com/alexfalkowski/go-service/v2/database/sql.DBs].
//
// Start with [Config], [Register], and [Open].
package pg

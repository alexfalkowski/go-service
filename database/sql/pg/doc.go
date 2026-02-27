// Package pg provides PostgreSQL (`database/sql`) wiring and helpers for go-service.
//
// This package integrates the pgx PostgreSQL driver with go-service SQL wiring by:
//
//   - registering the pgx stdlib `database/sql` driver under the name "pg" with OpenTelemetry instrumentation, and
//   - providing an `Open` constructor that creates master/slave connection pools using the shared SQL driver helpers.
//
// # Configuration and enablement
//
// PostgreSQL configuration is optional. By convention, a nil `*pg.Config` (or nil embedded config) is treated as
// "disabled", and constructors such as `Open` return (nil, nil) when disabled.
//
// # Master/slave pools
//
// `Open` resolves master and replica DSNs from configuration (DSNs are expressed as go-service "source strings"),
// connects using a master/slave pool abstraction, applies pool settings (max lifetime/open/idle), and registers
// OpenTelemetry DB stats metrics.
//
// Start with `Config`, `Register`, and `Open`.
package pg

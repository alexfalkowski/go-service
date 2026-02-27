// Package sql provides SQL database wiring and helpers for go-service.
//
// This package is the entrypoint for wiring SQL database support into Fx/Dig. It composes
// driver-specific modules (for example PostgreSQL via `database/sql/pg`) and exposes a small,
// consistent configuration shape for services (`sql.Config`).
//
// # Configuration and enablement
//
// SQL configuration is optional. By convention across go-service config types, a nil `*sql.Config`
// is treated as "disabled". Driver-specific configs typically follow the same convention, and
// constructors that depend on config often return (nil, nil) when disabled.
//
// # Master/slave pools and telemetry
//
// Driver integrations typically open master/slave connection pools, configure pool limits/lifetimes,
// and register OpenTelemetry instrumentation/metrics for `database/sql` connections.
//
// Start with `Module` and `Config`.
package sql

// Package sql provides SQL database wiring and helpers for go-service.
//
// This package is the entrypoint for wiring SQL database support into [go.uber.org/fx]/[go.uber.org/dig]. It composes
// driver-specific modules (for example PostgreSQL via [github.com/alexfalkowski/go-service/v2/database/sql/pg]) and exposes a small,
// consistent configuration shape for services ([Config]).
//
// # Configuration and enablement
//
// SQL configuration is optional. By convention across go-service config types, a nil *[Config]
// is treated as "disabled". Driver-specific configs typically follow the same convention, and
// constructors that depend on config often return (nil, nil) when disabled.
//
// # Writer/reader pool wrapper
//
// This package also exposes the shared writer/reader pool abstraction used by
// repository code:
//   - [DBs], the go-service writer/reader pool collection type, and
//   - [ConnectWritersReaders], the go-service helper for opening those pools.
//
// Callers choose the pool role explicitly with [DBs.Reader] or [DBs.Writer] and
// then use standard-library [database/sql] methods on the returned pool.
//
// This keeps internal code on a go-service import path while go-service-owned
// cleanup remains attached to [DBs.Destroy].
//
// # Writer/reader pools and telemetry
//
// Driver integrations typically open writer/reader connection pools, configure
// role-specific pool limits/lifetimes, wrap drivers with OpenTelemetry
// instrumentation when tracing or metrics are enabled, and register
// [database/sql] stats metrics when metrics are enabled.
//
// Start with [Module] and [Config].
package sql

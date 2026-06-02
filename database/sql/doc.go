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
// # Master/slave pool wrapper
//
// This package also exposes the shared master/slave pool abstraction used by
// repository code:
//   - [DBs], the go-service wrapper for the upstream pool collection type, and
//   - [ConnectMasterSlaves], the go-service wrapper for opening those pools.
//
// This keeps internal code on a go-service import path instead of importing the
// upstream helper package directly where the package graph allows it. The
// wrapper embeds the upstream type so existing query/ping/pool helper methods
// are still available, while go-service-owned cleanup remains attached to
// [DBs.Destroy].
//
// # Master/slave pools and telemetry
//
// Driver integrations typically open master/slave connection pools, configure pool limits/lifetimes,
// and register OpenTelemetry instrumentation/metrics for [database/sql] connections.
//
// Start with [Module] and [Config].
package sql

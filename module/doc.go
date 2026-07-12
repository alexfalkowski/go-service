// Package module provides top-level Fx module composition for go-service.
//
// This package defines opinionated, high-level module bundles that compose multiple lower-level
// feature modules into a single `di.Option` suitable for inclusion in an Fx/Dig application graph.
// These bundles are the primary supported entrypoints for applications and are the defaults used by
// `go-service-template` for long-running servers and `go-client-template` for short-lived commands.
//
// # Bundles
//
// The exported bundles are intended as defaults:
//
//   - [Library]: shared, transport-agnostic foundation wiring. It does not decode service
//     configuration by itself.
//
//   - [Server]: a typical long-running server composition, including configuration, telemetry,
//     transports, debug, and health wiring.
//
//   - [Client]: a typical short-lived or batch/client composition, including configuration,
//     telemetry, SQL/cache helpers, shared webhook helpers, and limiter key wiring.
//
// # Enablement and configuration
//
// In [Server] and [Client], most runtime subsystems are enabled/disabled by configuration, typically
// using optional pointer sub-configs (nil meaning "disabled"). These bundles wire constructors and
// registrations; whether a subsystem is active depends on the configuration supplied to the graph.
//
// Start with [Library], [Server], and [Client]. Drop down to lower-level package composition only when
// you intentionally need custom wiring beyond what the standard bundles provide.
package module

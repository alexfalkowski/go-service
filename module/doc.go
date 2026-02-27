// Package module provides top-level Fx module composition for go-service.
//
// This package defines opinionated, high-level module bundles that compose multiple lower-level
// feature modules into a single `di.Option` suitable for inclusion in an Fx/Dig application graph.
//
// # Bundles
//
// The exported bundles are intended as defaults:
//
//   - `Library`: shared, transport-agnostic wiring (building block for both servers and clients).
//
//   - `Server`: a typical server composition (builds on Library and adds configuration decoding,
//     server-side transports, telemetry, debugging, health checks, and common integrations).
//
//   - `Client`: a typical client composition (builds on Library and adds configuration decoding,
//     client-side transports/integrations, telemetry, and common client helpers).
//
// # Enablement and configuration
//
// Most subsystems in go-service are enabled/disabled by configuration, typically using optional
// pointer sub-configs (nil meaning "disabled"). These bundles wire constructors and registrations;
// whether a subsystem is active depends on the configuration supplied to the graph.
//
// Start with `Library`, `Server`, and `Client`.
package module

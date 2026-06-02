// Package module provides top-level Fx module composition for go-service.
//
// This package defines opinionated, high-level module bundles that compose multiple lower-level
// feature modules into a single `di.Option` suitable for inclusion in an Fx/Dig application graph.
// These bundles are the primary supported entrypoints for service applications and are the defaults
// used by `go-service-template`.
//
// # Bundles
//
// The exported bundles are intended as defaults:
//
//   - [Library]: shared, transport-agnostic wiring.
//
//   - [Server]: a typical server composition.
//
//   - [Client]: a typical client composition.
//
// # Enablement and configuration
//
// Most subsystems in go-service are enabled/disabled by configuration, typically using optional
// pointer sub-configs (nil meaning "disabled"). These bundles wire constructors and registrations;
// whether a subsystem is active depends on the configuration supplied to the graph.
//
// Start with [Library], [Server], and [Client]. Drop down to lower-level package composition only when
// you intentionally need custom wiring beyond what the standard bundles provide.
package module

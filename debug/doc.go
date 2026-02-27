// Package debug provides debug server wiring and diagnostic endpoints for go-service.
//
// This package wires an optional HTTP debug server that can expose operational and profiling endpoints.
// It is intended to be enabled in development and operations environments and disabled otherwise.
//
// # What this package provides
//
// At a high level, the debug subsystem consists of:
//
//   - A debug HTTP server (`Server`) that runs independently from the main service server.
//   - A debug router (`debug/http.ServeMux`) used to register endpoints.
//   - Optional endpoint registrations provided by subpackages (e.g. pprof, fgprof, statsviz, psutil).
//
// The actual endpoints are registered via the debug module wiring (see `Module`) which composes the
// subpackages and installs their handlers onto the debug mux.
//
// # Configuration and enablement
//
// Debug configuration is optional. By convention across go-service config types, a nil `*debug.Config`
// (or nil embedded config) is treated as "disabled", and `NewServer` returns (nil, nil) when disabled.
//
// # TLS
//
// When TLS config is enabled for the debug server, the server loads the configured certificate and key
// material using go-service "source strings" and constructs a `*crypto/tls.Config` via `crypto/tls.NewConfig`.
//
// Start with `Module` and `NewServer`.
package debug

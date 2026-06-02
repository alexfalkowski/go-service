// Package debug provides debug server wiring and diagnostic endpoints for go-service.
//
// This package wires an optional HTTP debug server that can expose operational and profiling endpoints.
// It is intended to be enabled in development and operations environments and disabled otherwise.
//
// # What this package provides
//
// At a high level, the debug subsystem consists of:
//
//   - A debug HTTP server ([Server]) that runs independently from the main service server.
//   - A debug router (debug/http.ServeMux) used to register endpoints.
//   - Built-in endpoint registrations installed through [Register].
//
// [Register] is the public front door for endpoint registration. The debug module wiring
// uses it to install the built-in handlers onto the debug mux, and callers that construct
// a debug server manually should use it instead of depending on endpoint implementation packages.
//
// # Configuration and enablement
//
// Debug configuration is optional. By convention across go-service config types, a nil *[Config]
// (or nil embedded config) is treated as "disabled", and [NewServer] returns (nil, nil) when disabled.
//
// When debug is enabled and no address is configured, the server binds to tcp://:6060. Production
// deployments should set an explicit address, TLS/mTLS, and network or policy controls appropriate
// for the environment.
//
// # TLS
//
// When TLS is enabled for the debug server, the package uses
// config/server.NewConfig to resolve crypto/tls/config.Config source
// strings and build the runtime *crypto/tls.Config assigned to the underlying
// HTTP server.
//
// Start with [Module], [NewServer], and [Register].
package debug

// Package metrics provides HTTP transport metrics instrumentation wiring for go-service.
//
// This package integrates metrics collection into the HTTP transport stack.
//
// Start with `Module` and `Register`.
//
// Registration: some transports use package-level registration to inject filesystem access or instrumentation.
// If you enable features that require registration, call `Register` during application startup before constructing clients/servers.
package metrics

// Package transport contains shared transport wiring and helpers for services built with go-service.
//
// This package provides common transport concepts and DI wiring that are used by the HTTP and gRPC transport stacks.
//
// Start with `Module` and `Register`.
//
// Registration: some transports use package-level registration to inject filesystem access or instrumentation.
// If you enable features that require registration, call `Register` during application startup before constructing clients/servers.
package transport

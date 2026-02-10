// Package health provides gRPC health transport wiring for go-service.
//
// This package integrates health checks with the gRPC transport stack.
//
// Start with `Module` and `Register`.
//
// Registration: some transports use package-level registration to inject filesystem access or instrumentation.
// If you enable features that require registration, call `Register` during application startup before constructing clients/servers.
package health

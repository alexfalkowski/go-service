// Package health provides HTTP health transport wiring for go-service.
//
// This package integrates health checks with the HTTP transport stack, typically by registering health endpoints
// on an HTTP server.
//
// Start with `Module` and `Register`.
//
// Registration: some transports use package-level registration to inject filesystem access or instrumentation.
// If you enable features that require registration, call `Register` during application startup before constructing clients/servers.
package health

// Package health provides HTTP health transport wiring for go-service.
//
// This package integrates health checks with the HTTP transport stack, typically by registering health endpoints
// on an HTTP server.
//
// Start with [Module] and [Register].
//
// [Register] registers service-prefixed /healthz, /livez, and /readyz handlers on the configured mux.
// [Module] wires that route registration into the HTTP transport graph.
package health

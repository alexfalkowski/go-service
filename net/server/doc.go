// Package server provides transport-agnostic server lifecycle helpers used by go-service.
//
// This package defines small primitives that higher-level transport packages (for example HTTP and gRPC)
// can use to run servers consistently and integrate with process shutdown logic.
//
// # Concepts
//
// The main abstractions are:
//
//   - Server: a minimal interface describing a runnable server that can be gracefully shut down.
//
//   - Service: a small lifecycle manager that starts a Server asynchronously, logs start/stop events,
//     and triggers application shutdown when the underlying Server.Serve returns an error.
//
// # Typical usage
//
// Transport-specific packages construct a concrete implementation of Server (wrapping a net/http.Server
// or grpc.Server), then wrap it in a Service and call Service.Start during application startup.
// On shutdown, the application stops the underlying server by calling Service.Stop with a context.
//
// Start with `Server` and `Service`.
package server

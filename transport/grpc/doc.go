// Package grpc contains gRPC transport utilities and wiring for services built with go-service.
//
// This package primarily exposes an Fx module (`Module`) that composes the building blocks needed to run and
// instrument gRPC servers and clients.
//
// Start with `Module` and `Register`.
//
// Registration: some transports use package-level registration to inject filesystem access or instrumentation.
// If you enable features that require registration, call `Register` during application startup before constructing clients/servers.
package grpc

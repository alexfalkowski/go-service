// Package grpc contains gRPC transport wiring for services built with go-service.
//
// This package primarily exposes an Fx module (`Module`) that composes the building blocks needed to run and
// instrument gRPC servers.
//
// Registration: this package uses package-level registration to inject filesystem access used when constructing
// TLS configuration. If you enable TLS, ensure `Register` is invoked during application startup before constructing
// clients/servers so the filesystem dependency is available.
//
// Start with `Module` for server wiring and `NewClient` / `NewDialOptions` for client construction.
package grpc

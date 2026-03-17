// Package health provides gRPC health protocol wiring for go-service.
//
// It contains the standard gRPC health service implementation and Fx wiring used to register that
// service with a gRPC server.
//
// Start with `NewServer`, `Register`, and `Module`.
package health

// Package health provides go-service aliases for the standard gRPC health protocol.
//
// This package re-exports google.golang.org/grpc/health/grpc_health_v1 behind a
// go-service import path. It intentionally does not add transport wiring,
// application health checks, or lifecycle behavior.
//
// Use transport/grpc/health for the go-service adapter that exposes application
// health state through this protocol.
package health

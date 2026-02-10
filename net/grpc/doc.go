// Package grpc provides small gRPC wrappers and helpers around google.golang.org/grpc.
//
// This package primarily re-exports common gRPC types and options behind go-service aliases and provides a few
// convenience helpers (for example interceptor chaining and credentials helpers) that are used by transport wiring.
//
// For full client/server wiring, see the higher-level transport packages under `transport/grpc` and module composition
// under `module`.
package grpc

// Package grpc provides small, stable wrappers and aliases around
// google.golang.org/grpc for use within go-service.
//
// This package deliberately exposes a narrow surface area:
//
//   - Type aliases for commonly used gRPC types (for example CallOption,
//     DialOption, ServerOption, Server, ClientConn, interceptors, and stream
//     types). These allow go-service code to depend on go-service packages while
//     still using the underlying gRPC implementations.
//
//   - Thin helper functions that forward to gRPC constructors/options (for
//     example StatsHandler, WithStatsHandler, ChainUnaryInterceptor,
//     ChainStreamInterceptor, Creds, NewTLS, NewInsecureCredentials, and
//     UseCompressor).
//
//   - A convenience NewServer constructor that applies standard server-side
//     keepalive configuration and registers gRPC reflection.
//
// # Server construction
//
// NewServer builds a *grpc.Server with keepalive enforcement and server
// parameters. Configuration values are sourced from an options.Map using the
// following keys (each value is a duration):
//
//   - keepalive_enforcement_policy_ping_min_time
//   - keepalive_max_connection_idle
//   - keepalive_max_connection_age
//   - keepalive_max_connection_age_grace
//   - keepalive_ping_time
//
// The timeout argument is used as the default value for each key when it is not
// present in the options map, and is also used as the keepalive ping Timeout.
//
// NewServer always enables server reflection via reflection.Register.
//
// # Telemetry and higher-level wiring
//
// This package is not responsible for full transport wiring (listeners, dial
// targets, TLS material loading, interceptors, lifecycle management, etc.).
// Those concerns are handled by higher-level packages (for example the transport
// and module packages). For OpenTelemetry gRPC stats handlers, see the sibling
// package net/grpc/telemetry.
package grpc

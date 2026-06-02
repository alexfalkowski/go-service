// Package server provides helpers for running a gRPC server as a managed go-service server.
//
// This package is intentionally small: it adapts a configured
// [github.com/alexfalkowski/go-service/v2/net/grpc.Server] to the go-service `server.Service` lifecycle,
// and it provides a lightweight wrapper that binds a listener and exposes a
// [Serve]/[Shutdown] API compatible with the generic server runner.
//
// In most cases you will:
//
//  1. Construct a *grpc.Server (typically via net/grpc.NewServer).
//  2. Provide a bind address via [github.com/alexfalkowski/go-service/v2/net/grpc/config.Config].
//  3. Call NewService to get a `server.Service` that starts and stops the gRPC
//     server and wires shutdown and logging.
//
// # Address format
//
// The bind address is read from [github.com/alexfalkowski/go-service/v2/net/grpc/config.Config.Address]. It may use the
// go-service network address format (for example "tcp://:9090") or a raw listen
// address such as ":9090", which defaults to the "tcp" network. Internally this
// is split into network/address and then passed to [net.Listen].
//
// # Relationship to other packages
//
//   - net/grpc/config: defines the gRPC bind address config type.
//   - net/grpc: provides small gRPC aliases and helpers (keepalive, reflection, etc.).
//   - transport/grpc (elsewhere in the repo): typically provides higher-level
//     client/server wiring (TLS material, interceptors, telemetry, etc.).
package server

// Package grpc contains gRPC transport wiring for services built with go-service.
//
// It provides constructors, interceptors, and Fx module wiring for:
//   - gRPC servers (`NewServer`) with standardized interceptors (metadata, logging, auth, rate limiting, etc.).
//   - gRPC clients (`NewClient`) with standardized dial options and client interceptors (timeouts, retries, breakers, auth, etc.).
//
// The primary entrypoint for DI wiring is `Module`, which composes this package with supporting subpackages:
// breaker, retry, limiter, metadata extraction/injection, token auth, and health.
//
// # Server wiring
//
// `NewServer` constructs a `*Server` wrapper around a configured `*grpc.Server` and a `*server.Service` lifecycle
// helper. When the transport is disabled via config, constructors in this package typically return nil so that
// downstream wiring can treat the server as "not enabled".
//
// # Client wiring
//
// `NewDialOptions` builds a slice of `grpc.DialOption` based on provided `ClientOption`s. `NewClient` is a convenience
// that applies those options and dials a `*grpc.ClientConn` to a target.
//
// Client-side concerns are expressed via options such as `WithClientTimeout`, `WithClientRetry`, `WithClientBreaker`,
// `WithClientLimiter`, and token-generator options. These options configure which interceptors are installed.
//
// # Registration gotcha (TLS filesystem)
//
// This package uses package-level registration to inject filesystem access used when constructing TLS configuration.
// The registered filesystem is used to resolve certificate/key "source strings" (for example `file:/path/to/cert` or
// `env:VAR`) during TLS configuration.
//
// If you enable TLS and do not call `Register` before constructing clients/servers, TLS configuration may fail to
// load key material.
package grpc

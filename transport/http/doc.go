// Package http provides HTTP transport wiring for services built with go-service.
//
// It provides constructors and Fx module wiring for HTTP servers and HTTP clients with standardized
// concerns such as:
//
//   - content negotiation and request/response encoding (`net/http/content`)
//   - MVC view rendering support (`net/http/mvc`)
//   - RPC and REST routing helpers (`net/http/rpc`, `net/http/rest`)
//   - request metadata extraction and propagation (`net/http/meta`)
//   - token authentication (server-side verification and client-side injection) (`transport/http/token`)
//   - client retries and circuit breakers (`transport/http/retry`, `transport/http/breaker`)
//   - server-side and client-side rate limiting (`transport/http/limiter`)
//   - health endpoints wiring (`transport/http/health`)
//   - Prometheus metrics endpoint wiring (`transport/http/telemetry/metrics`)
//
// Common lower-level HTTP primitives from `net/http` are also re-exported here so transport consumers can
// stay on a single import path. This package still focuses on higher-level server/client composition,
// middleware policy, and Fx wiring.
//
// The primary entrypoint for DI consumers is `Module`, which composes the HTTP transport stack and
// registers handlers/constructors needed to run an HTTP server. In typical service applications this
// happens through `module.Server` and `go-service-template`, so most consumers never call lower-level
// registration helpers directly.
//
// # Server wiring
//
// `NewServer` constructs an HTTP server service when enabled via config. When HTTP is disabled in the
// transport config, constructors in this package typically return nil so downstream wiring can treat
// the server as "not enabled".
//
// The server middleware chain typically includes metadata extraction, optional logging, optional token
// verification, and optional rate limiting. Health and metrics paths are treated as ignorable by some
// middleware (for example, token verification and rate limiting), so they do not require auth and do
// not consume limiter capacity by default.
//
// # Client wiring
//
// `NewClient` constructs an `*http.Client` whose `Transport` is assembled by `NewRoundTripper`. The
// resulting client can include metadata propagation, optional logging, optional token injection, optional
// retries, optional circuit breaking, and optional client-side rate limiting, depending on the provided
// client options.
//
// # Manual composition note (TLS filesystem)
//
// This package uses package-level registration to inject filesystem access used when constructing TLS
// configuration. The registered filesystem is used by
// `config/server.NewConfig` and `config/client.NewConfig` to resolve TLS "source strings"
// (for example `file:/path/to/cert` or `env:VAR`) when materializing the
// runtime `*crypto/tls.Config`.
//
// When you use `Module` (directly or through higher-level bundles such as `module.Server`), DI performs
// this registration for you. Call `Register` yourself only when you intentionally compose the HTTP transport
// package manually outside the standard module graph.
package http

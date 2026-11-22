// Package http provides HTTP transport wiring for services built with go-service.
//
// It provides constructors and go-signal module wiring for HTTP servers and HTTP clients with standardized
// concerns such as:
//
//   - content negotiation and request/response encoding ([github.com/alexfalkowski/go-service/v3/net/http/content/unary])
//   - MVC view rendering support ([github.com/alexfalkowski/go-service/v3/net/http/mvc])
//   - RPC and REST routing helpers ([github.com/alexfalkowski/go-service/v3/net/http/rpc], [github.com/alexfalkowski/go-service/v3/net/http/rest])
//   - request metadata extraction and propagation ([github.com/alexfalkowski/go-service/v3/net/http/meta])
//   - token authentication (server-side verification and client-side injection) ([github.com/alexfalkowski/go-service/v3/transport/http/token])
//   - client retries and circuit breakers ([github.com/alexfalkowski/go-service/v3/transport/http/retry], [github.com/alexfalkowski/go-service/v3/transport/http/breaker])
//   - server-side and client-side rate limiting ([github.com/alexfalkowski/go-service/v3/transport/http/limiter])
//   - health endpoints wiring ([github.com/alexfalkowski/go-service/v3/transport/http/health])
//   - Prometheus metrics endpoint wiring ([github.com/alexfalkowski/go-service/v3/transport/http/telemetry/metrics])
//
// Common lower-level HTTP primitives from [github.com/alexfalkowski/go-service/v3/net/http] are also re-exported here so transport consumers can
// stay on a single import path. This package still focuses on higher-level server/client composition,
// middleware policy, and go-signal wiring.
//
// The primary entrypoint for DI consumers is [Module], which composes the HTTP transport stack and
// registers handlers/constructors needed to run an HTTP server. In typical service applications this
// happens through [github.com/alexfalkowski/go-service/v3/module.Server] and `go-service-template`, so most consumers never call lower-level
// registration helpers directly.
//
// # Server wiring
//
// [NewServer] constructs an HTTP server service when enabled via config. When HTTP is disabled in the
// transport config, constructors in this package typically return nil so downstream wiring can treat
// the server as "not enabled".
//
// The server middleware chain typically includes metadata extraction, optional logging, optional token
// verification, and optional rate limiting. Route policy from [RoutePolicy] controls middleware bypasses:
// registered operation routes such as health and metrics do not require auth and do not consume limiter
// capacity by default, while registered unauthenticated routes bypass token verification and access control
// only.
// Non-streaming routes also receive the configured unary request timeout; streaming routes retain their
// stream-specific lifecycle controls.
//
// # Client wiring
//
// [NewClient] constructs an `*http.Client` using the middleware stack described by [NewRoundTripper].
// The stack can include metadata propagation, logging, token injection, retries, circuit breaking, and
// client-side rate limiting, depending on the provided options.
//
// [NewClientFromConfig] uses the configured endpoint, timeout, TLS, and breaker settings from
// [github.com/alexfalkowski/go-service/v3/config/client.Config]. It returns the validated base URL
// alongside the HTTP client so an SDK can use both. Supply retry policy through [WithClientRetry].
// To close idle connections at application shutdown, use [WithClientLifecycle] with either constructor.
// Construct the client before registering consumers so they stop before the pool closes.
//
// # Manual composition note (TLS filesystem)
//
// This package uses package-level registration to inject filesystem access used when constructing TLS
// configuration. The registered filesystem is used by
// [github.com/alexfalkowski/go-service/v3/config/server.NewConfig] and [github.com/alexfalkowski/go-service/v3/config/client.NewConfig] to resolve TLS "source strings"
// (for example `file:<path>` or `env:VAR`) when materializing the
// runtime `*crypto/tls.Config`.
//
// The standard [github.com/alexfalkowski/go-service/v3/module.Server] and
// [github.com/alexfalkowski/go-service/v3/module.Client] bundles register the filesystem before constructing OpenFeature
// providers. [Module] also registers it for direct HTTP transport composition. For manual setup, call
// [Register] before constructing a TLS client or server that reads file sources.
package http

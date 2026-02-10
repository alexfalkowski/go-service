// Package http provides the HTTP transport wiring for services built with this module.
//
// This package primarily exposes an Fx module (`Module`) that composes the building
// blocks needed to run and instrument HTTP servers and clients in go-service:
//
//   - mux and handler utilities from `github.com/alexfalkowski/go-service/v2/net/http`
//   - content helpers (`net/http/content`)
//   - MVC function map and registration (`net/http/mvc`)
//   - RPC and REST registration (`net/http/rpc`, `net/http/rest`)
//   - server-side limiting (`NewServerLimiter`)
//   - controller and token helpers (`NewController`, `NewToken`)
//   - token generation and verification (`transport/http/token`)
//   - HTTP metrics instrumentation (`transport/http/telemetry/metrics`)
//   - HTTP health transport integration (`transport/http/health`)
//
// Registration gotcha:
//
// TLS configuration may depend on a package-level filesystem registration pattern
// used by the transport stack. When enabling TLS, ensure the relevant transport
// registration has been performed (see `Register` in this package) before
// constructing servers or clients.
//
// The recommended entrypoint for DI consumers is `Module`.
package http

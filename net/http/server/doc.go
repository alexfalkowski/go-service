// Package server provides HTTP server adapters and lifecycle wiring for go-service.
//
// This package contains helpers used to run an internal HTTP server consistently within an Fx/Dig
// application. It bridges between:
//
//   - a configured `*net/http.Server` (from `github.com/alexfalkowski/go-service/v2/net/http`),
//   - a bind address (via `net.Listen`), and
//   - the transport-agnostic server lifecycle manager in `github.com/alexfalkowski/go-service/v2/net/server`.
//
// # Listener and address semantics
//
// `NewServer` parses `cfg.Address` using the go-service network address convention "<network>://<address>"
// (for example "tcp://:8080") and creates a listener using `net.Listen` with a background context.
// The resulting listener is held by the returned Server and is used for `Serve`/`ServeTLS`.
//
// # TLS semantics
//
// If `cfg.TLS` is non-nil, the Server assigns it to the underlying http.Server.TLSConfig and serves TLS
// using `ServeTLS` with empty cert/key paths (because certificates are expected to be provided via TLSConfig).
// If `cfg.TLS` is nil, the Server serves plain HTTP.
//
// # Error normalization
//
// `(*Server).Serve` wraps the underlying serve error using `net/http/errors.ServerError` so callers can
// treat expected shutdown errors consistently.
//
// # Service wiring
//
// `NewService` constructs a managed `*net/server.Service` that starts the HTTP server asynchronously,
// logs lifecycle events, and triggers application shutdown if the server terminates unexpectedly.
//
// Start with `NewService` and `NewServer`.
package server

// Package transport provides higher-level transport wiring for services built with go-service.
//
// It composes the HTTP and gRPC transport stacks into a single Fx module (`Module`) and provides lifecycle
// wiring via `net/server.Register` that starts and stops all configured server services.
//
// Lower-level protocol helpers and shared metadata utilities live under sibling `net/...` packages
// (for example `net/http/meta`, `net/grpc/meta`, `net/header`, and `net/server`).
//
// Typical usage is to include `transport.Module` in your application module graph and let DI construct the
// underlying transports:
//
// - `transport/http` for HTTP servers and clients.
// - `transport/grpc` for gRPC servers and clients.
// - `debug` (wired via `transport.NewServers`) for the optional debug server.
//
// # Registration gotchas
//
// Some transport subpackages require package-level registration to inject dependencies that cannot be
// automatically provided via constructors. In particular, the HTTP and gRPC transport packages use a
// registered filesystem when building TLS configuration (to read certificates/keys via source strings such
// as `file:` and `env:`). If you enable TLS and do not call the relevant transport registration prior to
// constructing clients/servers, TLS configuration may fail to load key material.
//
// When in doubt, call the transport-specific registration during application initialization (for example in
// your application's composition root) before constructing any transport servers or clients.
package transport

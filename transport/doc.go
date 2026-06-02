// Package transport provides higher-level transport wiring for services built with go-service.
//
// It composes the HTTP and gRPC transport stacks into a single Fx module ([Module]) and provides lifecycle
// wiring via [github.com/alexfalkowski/go-service/v2/net/server.Register] that starts and stops all configured server services.
//
// Lower-level protocol helpers and shared metadata utilities live under sibling net/... packages
// (for example [github.com/alexfalkowski/go-service/v2/net/http/meta],
// [github.com/alexfalkowski/go-service/v2/net/grpc/meta],
// [github.com/alexfalkowski/go-service/v2/net/header], and [github.com/alexfalkowski/go-service/v2/net/server]).
//
// Typical usage is to include [Module] in your application module graph and let DI construct the
// underlying transports. Services created from `go-service-template` usually reach this package indirectly
// via [github.com/alexfalkowski/go-service/v2/module.Server], so transport registration and lifecycle wiring are handled by the standard module
// graph.
//
// - [github.com/alexfalkowski/go-service/v2/transport/http] for HTTP servers and clients.
// - [github.com/alexfalkowski/go-service/v2/transport/grpc] for gRPC servers and clients.
// - debug (wired via [NewServers]) for the optional debug server.
//
// # Manual composition notes
//
// Some transport subpackages require package-level registration to inject dependencies that cannot be
// automatically provided via constructors. In particular, the HTTP and gRPC transport packages use a
// registered filesystem when building TLS configuration (to read certificates/keys via source strings such
// as `file:` and `env:`). The standard module graph performs this registration for you.
//
// Call the transport-specific registration yourself only when you intentionally bypass [Module]
// (or higher-level bundles such as [github.com/alexfalkowski/go-service/v2/module.Server]) and construct transport clients/servers manually.
package transport

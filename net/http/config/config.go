package config

import "crypto/tls"

// Config configures the internal HTTP server wiring.
//
// This config is used by go-service HTTP server adapters (for example `net/http/server.NewServer`) to
// create a listener and serve HTTP with optional TLS.
//
// It is intentionally minimal: it models only the bind address and an optional TLS configuration.
// Higher-level transport packages typically layer additional config (timeouts, protocol settings,
// middleware, etc.) elsewhere.
type Config struct {
	// TLS configures the TLS settings used by the HTTP server.
	//
	// When TLS is non-nil, server wiring typically assigns it to the underlying `net/http.Server.TLSConfig`
	// and serves TLS using `ServeTLS` with empty certificate/key paths (because the certificate material
	// is expected to be provided by the TLSConfig).
	//
	// When TLS is nil, the server is served over plain HTTP.
	TLS *tls.Config

	// Address is the bind address for the HTTP server.
	//
	// go-service commonly uses the "<network>://<address>" convention (for example "tcp://:8080") so the
	// network and address can be parsed and passed to net.Listen. Some callers may also supply a raw
	// host:port address depending on the surrounding wiring.
	Address string
}

package config

import "crypto/tls"

// Config configures the internal HTTP server wiring.
type Config struct {
	// TLS configures the server-side TLS settings used by the internal HTTP server.
	TLS *tls.Config

	// Address is the bind address for the internal HTTP server (for example ":8080").
	Address string
}

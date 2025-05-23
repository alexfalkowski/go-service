package config

import "crypto/tls"

// Config for HTTP.
type Config struct {
	TLS     *tls.Config
	Address string
}

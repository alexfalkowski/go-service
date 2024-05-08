package client

import (
	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/retry"
)

// IsEnabled for server.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for client.
type Config struct {
	Host      string        `yaml:"host,omitempty" json:"host,omitempty" toml:"host,omitempty"`
	Timeout   string        `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`
	TLS       *tls.Config   `yaml:"tls,omitempty" json:"tls,omitempty" toml:"tls,omitempty"`
	Retry     *retry.Config `yaml:"retry,omitempty" json:"retry,omitempty" toml:"retry,omitempty"`
	UserAgent string        `yaml:"user_agent,omitempty" json:"user_agent,omitempty" toml:"user_agent,omitempty"`
}

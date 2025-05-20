package server

import (
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/retry"
)

// IsEnabled for server.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for server.
type Config struct {
	Retry   *retry.Config `yaml:"retry,omitempty" json:"retry,omitempty" toml:"retry,omitempty"`
	TLS     *tls.Config   `yaml:"tls,omitempty" json:"tls,omitempty" toml:"tls,omitempty"`
	Timeout string        `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`
	Address string        `yaml:"address,omitempty" json:"address,omitempty" toml:"address,omitempty"`
}

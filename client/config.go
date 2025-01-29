package client

import (
	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/retry"
	"github.com/alexfalkowski/go-service/structs"
)

// IsEnabled for server.
func IsEnabled(cfg *Config) bool {
	return !structs.IsZero(cfg)
}

// Config for client.
type Config struct {
	TLS     *tls.Config   `yaml:"tls,omitempty" json:"tls,omitempty" toml:"tls,omitempty"`
	Retry   *retry.Config `yaml:"retry,omitempty" json:"retry,omitempty" toml:"retry,omitempty"`
	Address string        `yaml:"address,omitempty" json:"address,omitempty" toml:"address,omitempty"`
	Timeout string        `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`
}

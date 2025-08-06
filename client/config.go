package client

import (
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/retry"
)

// Config for client.
type Config struct {
	TLS     *tls.Config   `yaml:"tls,omitempty" json:"tls,omitempty" toml:"tls,omitempty"`
	Retry   *retry.Config `yaml:"retry,omitempty" json:"retry,omitempty" toml:"retry,omitempty"`
	Address string        `yaml:"address,omitempty" json:"address,omitempty" toml:"address,omitempty"`
	Timeout string        `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`
}

// IsEnabled for client.
func (c *Config) IsEnabled() bool {
	return c != nil
}

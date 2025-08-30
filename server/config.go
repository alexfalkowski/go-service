package server

import (
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/retry"
)

// Config for server.
type Config struct {
	Limiter *limiter.Config `yaml:"limiter,omitempty" json:"limiter,omitempty" toml:"limiter,omitempty"`
	Retry   *retry.Config   `yaml:"retry,omitempty" json:"retry,omitempty" toml:"retry,omitempty"`
	TLS     *tls.Config     `yaml:"tls,omitempty" json:"tls,omitempty" toml:"tls,omitempty"`
	Timeout string          `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`
	Address string          `yaml:"address,omitempty" json:"address,omitempty" toml:"address,omitempty"`
}

// IsEnabled for server.
func (c *Config) IsEnabled() bool {
	return c != nil
}

package server

import (
	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/retry"
)

// IsEnabled for server.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// UserAgent for server.
func UserAgent(cfg *Config) string {
	if !IsEnabled(cfg) {
		return ""
	}

	return cfg.UserAgent
}

// Config for server.
type Config struct {
	Port      string        `yaml:"port,omitempty" json:"port,omitempty" toml:"port,omitempty"`
	Retry     *retry.Config `yaml:"retry,omitempty" json:"retry,omitempty" toml:"retry,omitempty"`
	TLS       *tls.Config   `yaml:"tls,omitempty" json:"tls,omitempty" toml:"tls,omitempty"`
	UserAgent string        `yaml:"user_agent,omitempty" json:"user_agent,omitempty" toml:"user_agent,omitempty"`
}

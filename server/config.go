package server

import (
	"github.com/alexfalkowski/go-service/retry"
	"github.com/alexfalkowski/go-service/security"
)

// IsEnabled for server.
func IsEnabled(c *Config) bool {
	return c != nil && c.Enabled
}

// UserAgent for server.
func UserAgent(c *Config) string {
	if IsEnabled(c) {
		return c.UserAgent
	}

	return ""
}

// Config for server.
type Config struct {
	Enabled   bool             `yaml:"enabled,omitempty" json:"enabled,omitempty" toml:"enabled,omitempty"`
	Security  *security.Config `yaml:"security,omitempty" json:"security,omitempty" toml:"security,omitempty"`
	Port      string           `yaml:"port,omitempty" json:"port,omitempty" toml:"port,omitempty"`
	Retry     *retry.Config    `yaml:"retry,omitempty" json:"retry,omitempty" toml:"retry,omitempty"`
	UserAgent string           `yaml:"user_agent,omitempty" json:"user_agent,omitempty" toml:"user_agent,omitempty"`
}

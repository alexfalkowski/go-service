package server

import (
	"github.com/alexfalkowski/go-service/retry"
	"github.com/alexfalkowski/go-service/security"
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
	Security  *security.Config `yaml:"security,omitempty" json:"security,omitempty" toml:"security,omitempty"`
	Port      string           `yaml:"port,omitempty" json:"port,omitempty" toml:"port,omitempty"`
	Retry     *retry.Config    `yaml:"retry,omitempty" json:"retry,omitempty" toml:"retry,omitempty"`
	UserAgent string           `yaml:"user_agent,omitempty" json:"user_agent,omitempty" toml:"user_agent,omitempty"`
}

package http

import (
	"github.com/alexfalkowski/go-service/server"
)

// IsEnabled for HTTP.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && server.IsEnabled(cfg.Config)
}

// UserAgent for HTTP.
func UserAgent(cfg *Config) string {
	if !IsEnabled(cfg) {
		return ""
	}

	return server.UserAgent(cfg.Config)
}

// Config for HTTP.
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
}

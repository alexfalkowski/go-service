package http

import (
	"github.com/alexfalkowski/go-service/server"
)

// IsEnabled for HTTP.
func IsEnabled(c *Config) bool {
	return c != nil && server.IsEnabled(c.Config)
}

// UserAgent for HTTP.
func UserAgent(c *Config) string {
	if !IsEnabled(c) {
		return ""
	}

	return server.UserAgent(c.Config)
}

// Config for HTTP.
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
}

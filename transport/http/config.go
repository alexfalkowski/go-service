package http

import (
	"github.com/alexfalkowski/go-service/server"
)

// IsEnabled for http.
func IsEnabled(c *Config) bool {
	return c != nil && server.IsEnabled(c.Config)
}

// Config for HTTP.
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
}

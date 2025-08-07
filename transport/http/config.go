package http

import "github.com/alexfalkowski/go-service/v2/server"

// Config for HTTP.
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
}

// IsEnabled for HTTP.
func (c *Config) IsEnabled() bool {
	return c != nil && c.Config.IsEnabled()
}

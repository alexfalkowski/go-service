package debug

import "github.com/alexfalkowski/go-service/v2/server"

// Config for debug.
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
}

// IsEnabled for debug.
func (c *Config) IsEnabled() bool {
	return c != nil && c.Config.IsEnabled()
}

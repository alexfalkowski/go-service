package feature

import "github.com/alexfalkowski/go-service/v2/client"

// Config for feature.
type Config struct {
	*client.Config `yaml:",inline" json:",inline" toml:",inline"`
}

// IsEnabled for feature.
func (c *Config) IsEnabled() bool {
	return c != nil && c.Config.IsEnabled()
}

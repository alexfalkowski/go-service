package feature

import "github.com/alexfalkowski/go-service/v2/config/client"

// Config configures OpenFeature client behavior.
type Config struct {
	*client.Config `yaml:",inline" json:",inline" toml:",inline"`
}

// IsEnabled reports whether feature configuration is present and enabled.
func (c *Config) IsEnabled() bool {
	return c != nil && c.Config.IsEnabled()
}

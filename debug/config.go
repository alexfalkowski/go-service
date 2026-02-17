package debug

import "github.com/alexfalkowski/go-service/v2/config/server"

// Config configures the debug server.
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
}

// IsEnabled reports whether debug configuration is present and enabled.
func (c *Config) IsEnabled() bool {
	return c != nil && c.Config.IsEnabled()
}

package pg

import "github.com/alexfalkowski/go-service/v2/database/sql/config"

// Config contains PostgreSQL SQL database configuration.
//
// It embeds `database/sql/config`.Config to reuse common pool settings and DSN
// configuration.
type Config struct {
	*config.Config `yaml:",inline" json:",inline" toml:",inline"`
}

// IsEnabled reports whether PostgreSQL configuration is present and enabled.
func (c *Config) IsEnabled() bool {
	return c != nil && c.Config.IsEnabled()
}

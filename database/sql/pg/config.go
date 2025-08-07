package pg

import "github.com/alexfalkowski/go-service/v2/database/sql/config"

// Config for pg.
type Config struct {
	*config.Config `yaml:",inline" json:",inline" toml:",inline"`
}

// IsEnabled for pg.
func (c *Config) IsEnabled() bool {
	return c != nil && c.Config.IsEnabled()
}

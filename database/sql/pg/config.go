package pg

import (
	"github.com/alexfalkowski/go-service/database/sql/config"
)

// IsEnabled for pg.
func IsEnabled(c *Config) bool {
	return c != nil && config.IsEnabled(c.Config)
}

// Config for pg.
type Config struct {
	*config.Config `yaml:",inline" json:",inline" toml:",inline"`
}

package pg

import "github.com/alexfalkowski/go-service/v2/database/sql/config"

// IsEnabled for pg.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && config.IsEnabled(cfg.Config)
}

// Config for pg.
type Config struct {
	*config.Config `yaml:",inline" json:",inline" toml:",inline"`
}

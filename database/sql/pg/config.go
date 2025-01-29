package pg

import (
	"github.com/alexfalkowski/go-service/database/sql/config"
	"github.com/alexfalkowski/go-service/structs"
)

// IsEnabled for pg.
func IsEnabled(cfg *Config) bool {
	return !structs.IsZero(cfg)
}

// Config for pg.
type Config struct {
	*config.Config `yaml:",inline" json:",inline" toml:",inline"`
}

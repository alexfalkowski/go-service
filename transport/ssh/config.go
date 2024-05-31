package ssh

import (
	"github.com/alexfalkowski/go-service/server"
)

// IsEnabled for SSH.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && server.IsEnabled(cfg.Config)
}

// Config for SSH.
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
}

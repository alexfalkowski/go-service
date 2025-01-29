package proxy

import (
	"github.com/alexfalkowski/go-service/server"
)

// IsEnabled for proxy.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && server.IsEnabled(cfg.Config)
}

// Config for proxy.
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
}

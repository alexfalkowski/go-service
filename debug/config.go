package debug

import "github.com/alexfalkowski/go-service/v2/server"

// IsEnabled for debug.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && server.IsEnabled(cfg.Config)
}

// Config for debug.
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
}

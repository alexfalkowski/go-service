package feature

import "github.com/alexfalkowski/go-service/v2/client"

// IsEnabled for feature.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && client.IsEnabled(cfg.Config)
}

// Config for feature.
type Config struct {
	*client.Config `yaml:",inline" json:",inline" toml:",inline"`
}

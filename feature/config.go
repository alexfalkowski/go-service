package feature

import (
	"github.com/alexfalkowski/go-service/client"
)

// IsEnabled for feature.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && cfg.Kind != "" && client.IsEnabled(cfg.Config)
}

// Config for feature.
type Config struct {
	*client.Config `yaml:",inline" json:",inline" toml:",inline"`
	Kind           string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
}

// IsFlipt configuration.
func (c *Config) IsFlipt() bool {
	return c.Kind == "flipt"
}

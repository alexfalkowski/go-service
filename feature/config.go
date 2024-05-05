package feature

import (
	"github.com/alexfalkowski/go-service/client"
)

// IsEnabled for feature.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && cfg.Config != nil && cfg.Kind != ""
}

// Config for feature.
type Config struct {
	Kind           string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	*client.Config `yaml:",inline" json:",inline" toml:",inline"`
}

// IsFlipt configuration.
func (c *Config) IsFlipt() bool {
	return c.Kind == "flipt"
}

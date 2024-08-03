package tracer

import (
	"github.com/alexfalkowski/go-service/os"
)

// IsEnabled for tracer.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && cfg.Kind != ""
}

// Config for tracer.
type Config struct {
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	URL  string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty"`
	Key  string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
}

// GetKey for tracer.
func (c *Config) GetKey() (string, error) {
	return os.ReadFile(c.Key)
}

// HasKey for tracer.
func (c *Config) HasKey() bool {
	return c.Key != ""
}

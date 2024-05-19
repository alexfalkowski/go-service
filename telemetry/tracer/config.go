package tracer

import (
	"os"
	"path/filepath"
	"strings"
)

// IsEnabled for tracer.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && cfg.Kind != ""
}

type (
	// Key for tracer.
	Key string

	// Config for tracer.
	Config struct {
		Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
		Host string `yaml:"host,omitempty" json:"host,omitempty" toml:"host,omitempty"`
		Key  Key    `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
	}
)

// GetKey for tracer.
func (c *Config) GetKey() (string, error) {
	k, err := os.ReadFile(filepath.Clean(string(c.Key)))

	return strings.TrimSpace(string(k)), err
}

// HasKey for tracer.
func (c *Config) HasKey() bool {
	return c.Key != ""
}

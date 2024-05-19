package metrics

import (
	"os"
	"path/filepath"
	"strings"
)

// IsEnabled for metrics.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && cfg.Kind != ""
}

type (
	// Key for metrics.
	Key string

	// Config for metrics.
	Config struct {
		Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
		Host string `yaml:"host,omitempty" json:"host,omitempty" toml:"host,omitempty"`
		Key  Key    `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
	}
)

// IsOTLP configuration.
func (c *Config) IsOTLP() bool {
	return c.Kind == "otlp"
}

// IsPrometheus configuration.
func (c *Config) IsPrometheus() bool {
	return c.Kind == "prometheus"
}

// GetKey for metrics.
func (c *Config) GetKey() (string, error) {
	k, err := os.ReadFile(filepath.Clean(string(c.Key)))

	return strings.TrimSpace(string(k)), err
}

// HasKey for metrics.
func (c *Config) HasKey() bool {
	return c.Key != ""
}

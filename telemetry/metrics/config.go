package metrics

import (
	"github.com/alexfalkowski/go-service/os"
)

// IsEnabled for metrics.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && cfg.Kind != ""
}

// Config for metrics.
type Config struct {
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	Host string `yaml:"host,omitempty" json:"host,omitempty" toml:"host,omitempty"`
	Key  string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
}

// IsOTLP configuration.
func (c *Config) IsOTLP() bool {
	return c.Kind == "otlp"
}

// IsPrometheus configuration.
func (c *Config) IsPrometheus() bool {
	return c.Kind == "prometheus"
}

// GetKey for metrics.
func (c *Config) GetKey() string {
	return os.GetFromEnv(c.Key)
}

// HasKey for metrics.
func (c *Config) HasKey() bool {
	return c.Key != ""
}

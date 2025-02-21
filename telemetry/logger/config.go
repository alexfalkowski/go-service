package logger

import (
	"github.com/alexfalkowski/go-service/telemetry/header"
)

// IsEnabled for logger.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && cfg.Kind != ""
}

// Config for logger.
type Config struct {
	Headers header.Map `yaml:"headers,omitempty" json:"headers,omitempty" toml:"headers,omitempty"`
	Kind    string     `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	URL     string     `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty"`
	Level   string     `yaml:"level,omitempty" json:"level,omitempty" toml:"level,omitempty"`
}

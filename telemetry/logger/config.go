package logger

import (
	"github.com/alexfalkowski/go-service/strings"
	"github.com/alexfalkowski/go-service/telemetry/header"
)

// IsEnabled for logger.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && !strings.IsEmpty(cfg.Kind)
}

// Config for logger.
type Config struct {
	Headers header.Map `yaml:"headers,omitempty" json:"headers,omitempty" toml:"headers,omitempty"`
	Kind    string     `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	URL     string     `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty" validate:"omitempty,http_url"`
	Level   string     `yaml:"level,omitempty" json:"level,omitempty" toml:"level,omitempty"`
}

// IsOTLP configuration.
func (c *Config) IsOTLP() bool {
	return c.Kind == "otlp"
}

// IsJSON configuration.
func (c *Config) IsJSON() bool {
	return c.Kind == "json"
}

// IsText configuration.
func (c *Config) IsText() bool {
	return c.Kind == "text"
}

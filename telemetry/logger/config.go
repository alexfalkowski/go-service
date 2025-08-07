package logger

import "github.com/alexfalkowski/go-service/v2/telemetry/header"

// Config for logger.
type Config struct {
	Headers header.Map `yaml:"headers,omitempty" json:"headers,omitempty" toml:"headers,omitempty"`
	Kind    string     `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	URL     string     `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty" validate:"omitempty,http_url"`
	Level   string     `yaml:"level,omitempty" json:"level,omitempty" toml:"level,omitempty"`
}

// IsEnabled for logger.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// IsOTLP logger.
func (c *Config) IsOTLP() bool {
	return c.Kind == "otlp"
}

// IsJSON logger.
func (c *Config) IsJSON() bool {
	return c.Kind == "json"
}

// IsText logger.
func (c *Config) IsText() bool {
	return c.Kind == "text"
}

// IsTint logger.
func (c *Config) IsTint() bool {
	return c.Kind == "tint"
}

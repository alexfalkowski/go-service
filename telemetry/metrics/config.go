package metrics

import (
	"github.com/alexfalkowski/go-service/strings"
	"github.com/alexfalkowski/go-service/telemetry/header"
)

// IsEnabled for metrics.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && !strings.IsEmpty(cfg.Kind)
}

// Config for metrics.
type Config struct {
	Headers header.Map `yaml:"headers,omitempty" json:"headers,omitempty" toml:"headers,omitempty"`
	Kind    string     `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	URL     string     `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty" validate:"omitempty,http_url"`
}

// IsOTLP configuration.
func (c *Config) IsOTLP() bool {
	return c.Kind == "otlp"
}

// IsPrometheus configuration.
func (c *Config) IsPrometheus() bool {
	return c.Kind == "prometheus"
}

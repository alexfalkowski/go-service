package metrics

import "github.com/alexfalkowski/go-service/v2/telemetry/header"

// Config for metrics.
type Config struct {
	Headers header.Map `yaml:"headers,omitempty" json:"headers,omitempty" toml:"headers,omitempty"`
	Kind    string     `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	URL     string     `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty" validate:"omitempty,http_url"`
}

// IsEnabled for metrics.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// IsOTLP configuration.
func (c *Config) IsOTLP() bool {
	return c.Kind == "otlp"
}

// IsPrometheus configuration.
func (c *Config) IsPrometheus() bool {
	return c.Kind == "prometheus"
}

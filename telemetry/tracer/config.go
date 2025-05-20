package tracer

import (
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/header"
)

// IsEnabled for tracer.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && !strings.IsEmpty(cfg.Kind)
}

// Config for tracer.
type Config struct {
	Headers header.Map `yaml:"headers,omitempty" json:"headers,omitempty" toml:"headers,omitempty"`
	Kind    string     `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	URL     string     `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty" validate:"omitempty,http_url"`
}

package logger

import "github.com/alexfalkowski/go-service/v2/telemetry/header"

// Config configures telemetry logging.
type Config struct {
	// Headers contains exporter/request headers.
	//
	// Values may be configured as "source strings" (for example "env:NAME", "file:/path", or a literal value)
	// and are resolved by header.Map.Secrets / header.Map.MustSecrets.
	Headers header.Map `yaml:"headers,omitempty" json:"headers,omitempty" toml:"headers,omitempty"`

	// Kind selects the logger/exporter implementation (for example "stdout", "otlp", etc.),
	// depending on which implementations are compiled/registered by the service.
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`

	// URL is the destination endpoint for the selected Kind, when applicable.
	//
	// It must be a valid HTTP URL when set.
	URL string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty" validate:"omitempty,http_url"`

	// Level is the minimum log level to emit (for example "debug", "info", "warn", "error").
	Level string `yaml:"level,omitempty" json:"level,omitempty" toml:"level,omitempty"`
}

// IsEnabled for logger.
func (c *Config) IsEnabled() bool {
	return c != nil
}

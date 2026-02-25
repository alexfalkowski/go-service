package logger

import "github.com/alexfalkowski/go-service/v2/telemetry/header"

// Config configures telemetry logging.
type Config struct {
	// Headers contains exporter/request headers.
	//
	// Values may be configured as "source strings" (for example "env:NAME", "file:/path", or a literal value)
	// and are resolved by header.Map.Secrets / header.Map.MustSecrets.
	//
	// Headers are primarily used by exporters (for example OTLP) to pass authentication
	// or routing metadata to a collector.
	Headers header.Map `yaml:"headers,omitempty" json:"headers,omitempty" toml:"headers,omitempty"`

	// Kind selects the logger implementation.
	//
	// Supported kinds depend on what the service builds in, but this package typically
	// supports:
	//   - "otlp": export logs via OpenTelemetry OTLP/HTTP
	//   - "json": write JSON logs to stdout
	//   - "text": write text logs to stdout
	//   - "tint": write colorized text logs to stdout
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`

	// URL is the destination endpoint for the selected Kind, when applicable.
	//
	// For OTLP, this is the OTLP/HTTP logs endpoint URL. It must be a valid HTTP URL
	// when set.
	URL string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty" validate:"omitempty,http_url"`

	// Level is the minimum log level to emit.
	//
	// Values are interpreted as slog levels (for example "debug", "info", "warn", "error").
	// When unset, the behavior depends on the selected handler implementation.
	Level string `yaml:"level,omitempty" json:"level,omitempty" toml:"level,omitempty"`
}

// IsEnabled reports whether logging is enabled.
func (c *Config) IsEnabled() bool {
	return c != nil
}

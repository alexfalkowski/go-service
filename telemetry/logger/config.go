package logger

import "github.com/alexfalkowski/go-service/v2/telemetry/header"

// Config configures structured logging and optional log exporting.
//
// This package is responsible for constructing the logger based on Config.
// However, note that this type intentionally does not resolve header secrets by
// itself. If Headers contains go-service “source strings” (for example "env:NAME"
// or "file:/path"), callers are expected to resolve them before constructing
// exporters/handlers (for example via header.Map.Secrets or header.Map.MustSecrets).
//
// Enablement is modeled by presence: a nil *Config means logging is disabled.
type Config struct {
	// Headers contains exporter/request headers.
	//
	// These headers are primarily used by exporter-backed logger kinds (for example
	// "otlp") to pass authentication and/or routing metadata to a collector.
	//
	// Values may be configured as go-service “source strings” (for example "env:NAME",
	// "file:/path", or a literal value). Resolution is performed by the consumer that
	// prepares configuration for use by the logger/exporter.
	Headers header.Map `yaml:"headers,omitempty" json:"headers,omitempty" toml:"headers,omitempty"`

	// Kind selects the logger implementation.
	//
	// Supported kinds depend on what the service links in, but this package typically
	// supports:
	//
	//   - "otlp": export logs via OpenTelemetry OTLP/HTTP (and bridge slog records to OTel).
	//   - "json": write JSON logs to stdout.
	//   - "text": write text logs to stdout.
	//   - "tint": write colorized text logs to stdout.
	//
	// If Kind is unknown, logger construction returns ErrNotFound.
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`

	// URL is the destination endpoint for the selected Kind, when applicable.
	//
	// For "otlp", this is the OTLP/HTTP logs endpoint URL. It must be a valid HTTP URL
	// when set.
	URL string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty" validate:"omitempty,http_url"`

	// Level is the minimum log level to emit.
	//
	// Values are interpreted as slog levels (for example "debug", "info", "warn", "error").
	// When unset, the selected handler/exporter kind determines the effective default.
	Level string `yaml:"level,omitempty" json:"level,omitempty" toml:"level,omitempty"`
}

// IsEnabled reports whether logging configuration is present.
func (c *Config) IsEnabled() bool {
	return c != nil
}

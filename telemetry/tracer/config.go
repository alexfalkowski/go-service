package tracer

import "github.com/alexfalkowski/go-service/v2/telemetry/header"

// Config configures OpenTelemetry tracing export.
type Config struct {
	// Headers contains exporter/request headers.
	//
	// These headers are primarily used by exporter-backed tracer kinds (for example
	// "otlp") to pass authentication and/or routing metadata to a collector.
	//
	// Values may be configured using go-service “source strings” (for example "env:NAME",
	// "file:/path", or a literal value). This package does not resolve secrets itself;
	// resolution is performed by the consumer that prepares configuration for use by the
	// exporter (for example via header.Map.Secrets or header.Map.MustSecrets).
	Headers header.Map `yaml:"headers,omitempty" json:"headers,omitempty" toml:"headers,omitempty"`

	// Kind selects the tracer/exporter implementation.
	//
	// An empty kind means tracing is not configured. This package supports "otlp" and
	// wires an OTLP/HTTP exporter for that kind.
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`

	// URL is the destination endpoint for the selected Kind, when applicable.
	//
	// For "otlp", this is the OTLP/HTTP traces endpoint URL. It must be a valid HTTP URL
	// when set.
	URL string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty" validate:"omitempty,http_url"`
}

// IsEnabled reports whether tracing is configured.
//
// A nil *Config or empty Kind indicates tracing is disabled.
func (c *Config) IsEnabled() bool {
	return c != nil && c.Kind != ""
}

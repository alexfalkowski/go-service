package metrics

import "github.com/alexfalkowski/go-service/v2/telemetry/header"

// Config configures OpenTelemetry metrics export.
type Config struct {
	// Headers contains exporter/request headers.
	//
	// These headers are primarily used by exporter-based kinds (for example "otlp") to pass
	// authentication and/or routing metadata to a collector.
	//
	// Values may be configured as go-service “source strings” (for example "env:NAME",
	// "file:/path", or a literal value). Resolution is performed by the consumer that
	// prepares configuration for use by the exporter (for example via header.Map.Secrets
	// or header.Map.MustSecrets).
	Headers header.Map `yaml:"headers,omitempty" json:"headers,omitempty" toml:"headers,omitempty"`

	// Kind selects the metrics reader/exporter implementation.
	//
	// Supported kinds depend on what the service links in, but this package typically supports:
	//
	//   - "otlp": export metrics via OpenTelemetry OTLP/HTTP using a periodic reader.
	//   - "prometheus": expose metrics via the Prometheus exporter/reader.
	//
	// If Kind is unknown, reader construction will return an error (see the metrics package’s
	// ErrNotFound).
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`

	// URL is the destination endpoint for the selected Kind, when applicable.
	//
	// For "otlp", this is the OTLP/HTTP metrics endpoint URL. It must be a valid HTTP URL
	// when set.
	//
	// For "prometheus", URL is typically ignored by the exporter/reader implementation.
	URL string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty" validate:"omitempty,http_url"`
}

// IsEnabled reports whether metrics configuration is present.
//
// A nil *Config indicates metrics are disabled.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// IsPrometheus reports whether the configured Kind is "prometheus".
func (c *Config) IsPrometheus() bool {
	return c.Kind == "prometheus"
}

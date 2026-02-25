package metrics

import "github.com/alexfalkowski/go-service/v2/telemetry/header"

// Config configures OpenTelemetry metrics export.
type Config struct {
	// Headers contains exporter/request headers.
	//
	// Values may be configured as "source strings" (for example "env:NAME", "file:/path", or a literal value)
	// and are resolved by header.Map.Secrets / header.Map.MustSecrets.
	//
	// Headers are primarily used by exporters (for example OTLP) to pass authentication
	// or routing metadata to a collector.
	Headers header.Map `yaml:"headers,omitempty" json:"headers,omitempty" toml:"headers,omitempty"`

	// Kind selects the metrics reader/exporter implementation.
	//
	// Supported kinds depend on what the service builds in, but this package typically supports:
	//   - "otlp": export metrics via OpenTelemetry OTLP/HTTP using a periodic reader
	//   - "prometheus": expose metrics via the Prometheus exporter/reader
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`

	// URL is the destination endpoint for the selected Kind, when applicable.
	//
	// For OTLP, this is the OTLP/HTTP metrics endpoint URL. It must be a valid HTTP URL
	// when set.
	URL string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty" validate:"omitempty,http_url"`
}

// IsEnabled reports whether metrics configuration is enabled.
//
// A nil *Config indicates metrics are disabled.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// IsPrometheus reports whether the configured Kind is "prometheus".
func (c *Config) IsPrometheus() bool {
	return c.Kind == "prometheus"
}

package metrics

import "github.com/alexfalkowski/go-service/v2/telemetry/header"

// Config configures telemetry metrics.
type Config struct {
	// Headers contains exporter/request headers.
	//
	// Values may be configured as "source strings" (for example "env:NAME", "file:/path", or a literal value)
	// and are resolved by header.Map.Secrets / header.Map.MustSecrets.
	Headers header.Map `yaml:"headers,omitempty" json:"headers,omitempty" toml:"headers,omitempty"`

	// Kind selects the metrics exporter implementation (for example "prometheus", "otlp", etc.),
	// depending on which implementations are compiled/registered by the service.
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`

	// URL is the destination endpoint for the selected Kind, when applicable.
	//
	// It must be a valid HTTP URL when set.
	URL string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty" validate:"omitempty,http_url"`
}

// IsEnabled reports whether metrics configuration is enabled.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// IsPrometheus reports whether the configured exporter Kind is Prometheus.
func (c *Config) IsPrometheus() bool {
	return c.Kind == "prometheus"
}

package tracer

import "github.com/alexfalkowski/go-service/v2/telemetry/header"

// Config configures OpenTelemetry tracing export.
type Config struct {
	// Headers contains exporter/request headers.
	//
	// These headers are primarily used by exporter-backed tracer kinds (for example
	// "otlp") to pass authentication and/or routing metadata to a collector.
	//
	// Values may be configured using go-service "source strings" (for example "env:NAME",
	// "file:/path", or a literal value). This package does not resolve secrets itself;
	// resolution is performed by the consumer that prepares configuration for use by the
	// exporter (for example via header.Map.Secrets or header.Map.MustSecrets).
	Headers header.Map `yaml:"headers,omitempty" json:"headers,omitempty" toml:"headers,omitempty"`

	// Sampler configures trace head sampling.
	//
	// A nil or empty sampler preserves the OpenTelemetry SDK default sampler and
	// SDK environment handling. When set, it overrides those defaults.
	Sampler *SamplerConfig `yaml:"sampler,omitempty" json:"sampler,omitempty" toml:"sampler,omitempty"`

	// Kind selects the tracer/exporter implementation.
	//
	// An empty kind means tracing is not configured. This package supports "otlp" and
	// wires an OTLP/HTTP exporter for that kind.
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`

	// URL is the destination endpoint for the selected Kind, when applicable.
	//
	// For "otlp", this is the required OTLP/HTTP traces endpoint URL. It must be a
	// valid HTTP URL. Standard OpenTelemetry endpoint environment variables are not used as fallbacks;
	// configure this value explicitly through go-service config.
	URL string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty" validate:"omitempty,http_url"`
}

// IsEnabled reports whether tracing is configured.
//
// A nil *[Config] or empty Kind indicates tracing is disabled.
func (c *Config) IsEnabled() bool {
	return c != nil && c.Kind != ""
}

// SamplerConfig configures trace head sampling.
type SamplerConfig struct {
	// Kind selects the sampler implementation.
	//
	// Supported values are "always_on", "always_off", and "ratio".
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty" validate:"omitempty,oneof=always_on always_off ratio"`

	// Ratio is the fraction used by the ratio sampler when starting root traces.
	//
	// Values must be between 0 and 1, inclusive. A zero ratio drops new root
	// traces and a ratio of 1 samples every new root trace. Incoming parent
	// sampling decisions are preserved.
	Ratio float64 `yaml:"ratio,omitempty" json:"ratio,omitempty" toml:"ratio,omitempty" validate:"gte=0,lte=1"`
}

// IsEnabled reports whether sampler configuration is present.
func (c *SamplerConfig) IsEnabled() bool {
	return c != nil && c.Kind != ""
}

package metrics

import (
	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/telemetry/header"
	"github.com/alexfalkowski/go-service/v2/telemetry/internal/otlp"
	"github.com/alexfalkowski/go-service/v2/time"
)

// Config configures OpenTelemetry metrics export.
type Config struct {
	// Headers contains exporter/request headers.
	//
	// These headers are primarily used by exporter-based kinds (for example "otlp") to pass
	// authentication and/or routing metadata to a collector.
	//
	// Values may be configured as go-service "source strings" (for example "env:NAME",
	// "file:/path", or a literal value). Resolution is performed by the consumer that
	// prepares configuration for use by the exporter (for example via header.Map.Secrets
	// or header.Map.MustSecrets).
	Headers header.Map `yaml:"headers,omitempty" json:"headers,omitempty" toml:"headers,omitempty"`

	// TLS configures OTLP/gRPC client transport security.
	//
	// This only applies when Kind is "otlp" and Protocol is "grpc". A non-nil
	// config enables TLS; Cert, Key, and CA values use go-service source strings.
	TLS *tls.Config `yaml:"tls,omitempty" json:"tls,omitempty" toml:"tls,omitempty"`

	// Views maps OpenTelemetry instrument names to explicit histogram bucket
	// boundaries, overriding the SDK default boundaries for matching histogram
	// instruments such as "http.server.request.duration" or "rpc.server.call.duration".
	//
	// Instrument name matching follows OpenTelemetry view semantics and supports
	// "*" wildcards (for example "rpc.*.duration"). Boundaries are expressed in the
	// instrument's unit (seconds for duration histograms, bytes for size
	// histograms) and must be listed in increasing order. A nil or empty map keeps
	// the SDK default boundaries, so this field is a no-op unless configured.
	Views map[string][]float64 `yaml:"views,omitempty" json:"views,omitempty" toml:"views,omitempty"`

	// Kind selects the metrics reader/exporter implementation.
	//
	// Supported kinds depend on what the service links in, but this package typically supports:
	//
	//   - "otlp": export metrics via OpenTelemetry OTLP using a periodic reader.
	//   - "prometheus": expose metrics via the Prometheus exporter/reader.
	//
	// If Kind is unknown, reader construction will return an error (see the metrics package's
	// ErrNotFound).
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty" validate:"omitempty,oneof=otlp prometheus"`

	// Protocol selects the OTLP transport protocol.
	//
	// Supported values are "http" and "grpc". When empty, "http" is used.
	// This field only applies when Kind is "otlp".
	Protocol string `yaml:"protocol,omitempty" json:"protocol,omitempty" toml:"protocol,omitempty" validate:"omitempty,oneof=http grpc"`

	// URL is the destination endpoint for the selected Kind, when applicable.
	//
	// For "otlp" over HTTP, this is the required OTLP/HTTP metrics endpoint URL.
	// For "otlp" over gRPC, this is the required collector host:port endpoint.
	// Standard OpenTelemetry endpoint environment variables are not used as fallbacks;
	// configure this value explicitly through go-service config.
	//
	// For "prometheus", URL is typically ignored by the exporter/reader implementation.
	URL string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty"`

	// Interval is the OTLP periodic export interval.
	//
	// A zero value keeps the OpenTelemetry SDK default. Negative values are
	// invalid. This field only applies when Kind is "otlp".
	Interval time.Duration `yaml:"interval,omitempty" json:"interval,omitempty" toml:"interval,omitempty" validate:"gte=0"`

	// Timeout is the OTLP periodic export timeout.
	//
	// A zero value keeps the OpenTelemetry SDK default. Negative values are
	// invalid. This field only applies when Kind is "otlp".
	Timeout time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty" validate:"gte=0"`
}

// IsEnabled reports whether metrics configuration is present.
//
// A nil *[Config] indicates metrics are disabled.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// IsPrometheus reports whether the configured Kind is "prometheus".
func (c *Config) IsPrometheus() bool {
	return c.Kind == "prometheus"
}

// GetProtocol returns the configured OTLP transport protocol.
//
// A nil receiver or an empty value falls back to OTLP/HTTP.
func (c *Config) GetProtocol() string {
	if c == nil || c.Protocol == "" {
		return otlp.ProtocolHTTP
	}

	return c.Protocol
}

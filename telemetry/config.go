package telemetry

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
)

// Config configures service telemetry (logging, metrics, and tracing).
//
// This type acts as a single configuration root that services can embed into
// their overall configuration. Each field points at a per-signal configuration
// struct:
//
//   - Metadata configures context metadata copied to log records and spans.
//   - Logger configures application/system logging and any configured log exporters.
//   - Metrics configures metrics collection/readers/exporters.
//   - Propagation configures context propagation formats.
//   - Tracer configures distributed tracing (spans) and exporters.
//
// Enablement is intentionally modeled as presence: a nil *[Config] indicates that
// telemetry is disabled at the top level. Subpackages may also implement their
// own enable/disable semantics based on their specific config (for example nil
// config or an empty kind).
type Config struct {
	// Attributes are OpenTelemetry resource attributes attached to all configured
	// telemetry providers.
	//
	// Values are plain resource labels, not source strings. Fixed service identity
	// attributes such as service.name and service.version take precedence over
	// duplicate keys.
	Attributes attributes.Map `yaml:"attributes,omitempty" json:"attributes,omitempty" toml:"attributes,omitempty"`

	// Metadata configures request and service metadata exported to logs and traces.
	Metadata *MetadataConfig `yaml:"metadata,omitempty" json:"metadata,omitempty" toml:"metadata,omitempty"`

	// Logger configures application/system logging output and exporters.
	Logger *logger.Config `yaml:"logger,omitempty" json:"logger,omitempty" toml:"logger,omitempty"`

	// Metrics configures metrics collection and exporting.
	Metrics *metrics.Config `yaml:"metrics,omitempty" json:"metrics,omitempty" toml:"metrics,omitempty"`

	// Propagation configures OpenTelemetry context propagation.
	Propagation *PropagationConfig `yaml:"propagation,omitempty" json:"propagation,omitempty" toml:"propagation,omitempty"`

	// Tracer configures distributed tracing (spans) and exporting.
	Tracer *tracer.Config `yaml:"tracer,omitempty" json:"tracer,omitempty" toml:"tracer,omitempty"`
}

// MetadataConfig configures the metadata values exported to logs and traces.
//
// A zero MaxValueSize uses the 1,024-byte default. The original context
// metadata is not modified.
type MetadataConfig struct {
	// MaxValueSize caps each exported metadata value. It accepts the repository's
	// decimal byte-size syntax, such as "4KB".
	MaxValueSize bytes.Size `yaml:"max_value_size,omitempty" json:"max_value_size,omitempty" toml:"max_value_size,omitempty" validate:"config_size"`
}

// GetMaxValueSize returns the configured metadata value size limit.
//
// A nil receiver or zero value returns the 1,024-byte default.
func (c *MetadataConfig) GetMaxValueSize() meta.Limit {
	if c == nil || c.MaxValueSize == 0 {
		return meta.Limit(1024)
	}

	return meta.Limit(c.MaxValueSize)
}

// IsEnabled reports whether telemetry configuration is present.
//
// A nil receiver returns false, which is commonly used as a simple top-level
// enable/disable switch for telemetry wiring.
func (c *Config) IsEnabled() bool {
	return c != nil
}

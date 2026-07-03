package telemetry

import (
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

	// Logger configures application/system logging output and exporters.
	Logger *logger.Config `yaml:"logger,omitempty" json:"logger,omitempty" toml:"logger,omitempty"`

	// Metrics configures metrics collection and exporting.
	Metrics *metrics.Config `yaml:"metrics,omitempty" json:"metrics,omitempty" toml:"metrics,omitempty"`

	// Propagation configures OpenTelemetry context propagation.
	Propagation *PropagationConfig `yaml:"propagation,omitempty" json:"propagation,omitempty" toml:"propagation,omitempty"`

	// Tracer configures distributed tracing (spans) and exporting.
	Tracer *tracer.Config `yaml:"tracer,omitempty" json:"tracer,omitempty" toml:"tracer,omitempty"`
}

// IsEnabled reports whether telemetry configuration is present.
//
// A nil receiver returns false, which is commonly used as a simple top-level
// enable/disable switch for telemetry wiring.
func (c *Config) IsEnabled() bool {
	return c != nil
}

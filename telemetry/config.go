package telemetry

import (
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
//   - Tracer configures distributed tracing (spans) and exporters.
//
// Enablement is intentionally modeled as presence: a nil *Config indicates that
// telemetry is disabled at the top level. Subpackages may also implement their
// own enable/disable semantics based on their specific config (for example nil
// config or an empty kind).
type Config struct {
	// Logger configures application/system logging output and exporters.
	Logger *logger.Config `yaml:"logger,omitempty" json:"logger,omitempty" toml:"logger,omitempty"`

	// Metrics configures metrics collection and exporting.
	Metrics *metrics.Config `yaml:"metrics,omitempty" json:"metrics,omitempty" toml:"metrics,omitempty"`

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

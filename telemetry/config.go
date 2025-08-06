package telemetry

import (
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
)

// Config for telemetry.
type Config struct {
	Logger  *logger.Config  `yaml:"logger,omitempty" json:"logger,omitempty" toml:"logger,omitempty"`
	Metrics *metrics.Config `yaml:"metrics,omitempty" json:"metrics,omitempty" toml:"metrics,omitempty"`
	Tracer  *tracer.Config  `yaml:"tracer,omitempty" json:"tracer,omitempty" toml:"tracer,omitempty"`
}

// IsEnabled for telemetry.
func (c *Config) IsEnabled() bool {
	return c != nil
}

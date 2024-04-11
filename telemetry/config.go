package telemetry

import (
	"github.com/alexfalkowski/go-service/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
)

// Config for telemetry.
type Config struct {
	Logger  *zap.Config     `yaml:"logger,omitempty" json:"logger,omitempty" toml:"logger,omitempty"`
	Metrics *metrics.Config `yaml:"metrics,omitempty" json:"metrics,omitempty" toml:"metrics,omitempty"`
	Tracer  *tracer.Config  `yaml:"tracer,omitempty" json:"tracer,omitempty" toml:"tracer,omitempty"`
}

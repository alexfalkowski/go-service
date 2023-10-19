package telemetry

import (
	"github.com/alexfalkowski/go-service/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
)

// Config for telemetry.
type Config struct {
	Logger zap.Config    `yaml:"logger" json:"logger" toml:"logger"`
	Tracer tracer.Config `yaml:"tracer" json:"tracer" toml:"tracer"`
}

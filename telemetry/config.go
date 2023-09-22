package telemetry

import (
	"github.com/alexfalkowski/go-service/telemetry/tracer"
)

// Config for telemetry.
type Config struct {
	Tracer tracer.Config `yaml:"tracer" json:"tracer" toml:"tracer"`
}

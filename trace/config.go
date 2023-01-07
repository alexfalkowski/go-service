package trace

import (
	"github.com/alexfalkowski/go-service/trace/opentracing"
)

// Config for trace.
type Config struct {
	Opentracing opentracing.Config `yaml:"opentracing" json:"opentracing" toml:"opentracing"`
}

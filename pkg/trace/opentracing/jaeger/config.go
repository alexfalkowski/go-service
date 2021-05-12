package jaeger

import (
	"github.com/kelseyhightower/envconfig"
)

// NewConfig for jaeger.
func NewConfig() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)

	return &config, err
}

// Config for jaeger.
type Config struct {
	Name string `envconfig:"SERVICE_NAME" required:"true"`
	Host string `envconfig:"JAEGER_TRACE_HOST" required:"true" default:"localhost:6831"`
}

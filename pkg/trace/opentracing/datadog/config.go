package datadog

import (
	"github.com/kelseyhightower/envconfig"
)

// NewConfig for datadog.
func NewConfig() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)

	return &config, err
}

// Config for datadog.
type Config struct {
	Name string `envconfig:"SERVICE_NAME" required:"true"`
	Host string `envconfig:"DATADOG_TRACE_HOST" required:"true" default:"localhost:8126"`
}

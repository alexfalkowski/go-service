package http

import (
	"github.com/kelseyhightower/envconfig"
)

// NewConfig for HTTP.
func NewConfig() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)

	return &config, err
}

// Config for HTTP.
type Config struct {
	HTTPPort string `envconfig:"HTTP_PORT" required:"true" default:"8080"`
}

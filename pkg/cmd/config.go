package cmd

import (
	"github.com/kelseyhightower/envconfig"
)

// NewConfig for cmd.
func NewConfig() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)

	return &config, err
}

// Config for cmd.
type Config struct {
	Name        string `envconfig:"SERVICE_NAME" required:"true"`
	Description string `envconfig:"SERVICE_DESCRIPTION" required:"true"`
}

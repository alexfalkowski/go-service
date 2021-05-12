package pg

import (
	"github.com/kelseyhightower/envconfig"
)

// NewConfig for SQL.
func NewConfig() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)

	return &config, err
}

// Config for SQL.
type Config struct {
	Name string `envconfig:"SERVICE_NAME" required:"true"`
	URL  string `envconfig:"POSTGRESQL_URL" required:"true"`
}

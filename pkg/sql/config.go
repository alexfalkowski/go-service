package sql

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
	AppName     string `envconfig:"APP_NAME" required:"true"`
	PostgresURL string `envconfig:"POSTGRES_URL" required:"true"`
}

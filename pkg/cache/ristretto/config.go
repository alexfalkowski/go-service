package ristretto

import (
	"github.com/dgraph-io/ristretto"
	"github.com/kelseyhightower/envconfig"
)

// NewConfig for ristretto.
func NewConfig() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)

	config.Config = &ristretto.Config{
		NumCounters: 1e7,     // nolint:gomnd
		MaxCost:     1 << 30, // nolint:gomnd
		BufferItems: 64,      // nolint:gomnd
		Metrics:     true,
	}

	return &config, err
}

// Config for ristretto.
type Config struct {
	Name string `envconfig:"SERVICE_NAME" required:"true"`

	Config *ristretto.Config
}

package redis

import (
	"github.com/kelseyhightower/envconfig"
)

// NewConfig for redis.
func NewConfig() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)

	return &config, err
}

// Config for redis.
type Config struct {
	AppName string `envconfig:"APP_NAME" required:"true"`
	Host    string `envconfig:"REDIS_CACHE_HOST" required:"true" default:"localhost:6379"`
}

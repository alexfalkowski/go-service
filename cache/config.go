package cache

import (
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/structs"
)

// IsEnabled for cache.
func IsEnabled(cfg *Config) bool {
	return !structs.IsZero(cfg)
}

// Config for cache.
type Config struct {
	Redis *redis.Config `yaml:"redis,omitempty" json:"redis,omitempty" toml:"redis,omitempty"`
}

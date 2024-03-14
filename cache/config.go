package cache

import (
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/ristretto"
)

// Config for cache.
type Config struct {
	Redis     redis.Config     `yaml:"redis,omitempty" json:"redis,omitempty" toml:"redis,omitempty"`
	Ristretto ristretto.Config `yaml:"ristretto,omitempty" json:"ristretto,omitempty" toml:"ristretto,omitempty"`
}

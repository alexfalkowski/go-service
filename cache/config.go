package cache

import (
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/ristretto"
)

// Config for cache.
type Config struct {
	Redis     redis.Config     `yaml:"redis" json:"redis" toml:"redis"`
	Ristretto ristretto.Config `yaml:"ristretto" json:"ristretto" toml:"ristretto"`
}

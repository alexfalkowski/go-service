package cache

import (
	"github.com/alexfalkowski/go-service/pkg/cache/redis"
	"github.com/alexfalkowski/go-service/pkg/cache/ristretto"
)

// Config for cache.
type Config struct {
	Redis     redis.Config     `yaml:"redis"`
	Ristretto ristretto.Config `yaml:"ristretto"`
}

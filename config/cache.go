package config

import (
	"github.com/alexfalkowski/go-service/cache"
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/ristretto"
)

func (cfg *Config) RedisConfig() *redis.Config {
	if !cache.IsEnabled(cfg.Cache) {
		return nil
	}

	return cfg.Cache.Redis
}

func (cfg *Config) RistrettoConfig() *ristretto.Config {
	if !cache.IsEnabled(cfg.Cache) {
		return nil
	}

	return cfg.Cache.Ristretto
}

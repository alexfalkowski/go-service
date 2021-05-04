package cache

import (
	"github.com/alexfalkowski/go-service/pkg/cache/redis"
	"go.uber.org/fx"
)

var (
	// RedisModule for fx.
	RedisModule = fx.Provide(redis.NewCache)

	// Module for fx.
	Module = fx.Options(RedisModule)
)

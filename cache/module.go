package cache

import (
	"github.com/alexfalkowski/go-service/cache/redis"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	redis.Module,
)

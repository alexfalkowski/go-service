package cache

import (
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	redis.Module,
	ristretto.Module,
)

package redis

import (
	"github.com/alexfalkowski/go-service/cache/redis/telemetry/metrics"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewClient),
	fx.Provide(NewOptions),
	fx.Provide(NewCache),
	fx.Provide(NewRingOptions),
	fx.Invoke(metrics.Register),
)

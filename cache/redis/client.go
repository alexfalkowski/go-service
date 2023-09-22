package redis

import (
	"context"

	"github.com/alexfalkowski/go-service/cache/redis/client"
	"github.com/alexfalkowski/go-service/cache/redis/telemetry/logger"
	rzap "github.com/alexfalkowski/go-service/cache/redis/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/cache/redis/telemetry/tracer"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ClientParams for redis.
type ClientParams struct {
	fx.In

	Lifecycle   fx.Lifecycle
	RingOptions *redis.RingOptions
	Tracer      tracer.Tracer
	Logger      *zap.Logger
}

// NewClient for redis.
func NewClient(params ClientParams) client.Client {
	redis.SetLogger(logger.NewLogger())

	var client client.Client = redis.NewRing(params.RingOptions)
	client = tracer.NewClient(params.Tracer, client)
	client = rzap.NewClient(params.Logger, client)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client
}

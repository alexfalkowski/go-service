package redis

import (
	"context"

	"github.com/alexfalkowski/go-service/cache/redis/client"
	"github.com/alexfalkowski/go-service/cache/redis/logger"
	"github.com/alexfalkowski/go-service/cache/redis/trace/opentracing"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

// ClientParams for redis.
type ClientParams struct {
	fx.In

	Lifecycle   fx.Lifecycle
	RingOptions *redis.RingOptions
	Tracer      opentracing.Tracer
}

// NewClient for redis.
func NewClient(params ClientParams) client.Client {
	redis.SetLogger(logger.NewLogger())

	var client client.Client = redis.NewRing(params.RingOptions)
	client = opentracing.NewClient(params.Tracer, client)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client
}

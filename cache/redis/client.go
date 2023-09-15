package redis

import (
	"context"

	"github.com/alexfalkowski/go-service/cache/redis/client"
	"github.com/alexfalkowski/go-service/cache/redis/telemetry"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ClientParams for redis.
type ClientParams struct {
	fx.In

	Lifecycle   fx.Lifecycle
	RingOptions *redis.RingOptions
	Tracer      telemetry.Tracer
	Logger      *zap.Logger
}

// NewClient for redis.
func NewClient(params ClientParams) client.Client {
	redis.SetLogger(&logger{})

	var client client.Client = redis.NewRing(params.RingOptions)
	client = telemetry.NewTracerClient(params.Tracer, client)
	client = telemetry.NewLoggerClient(params.Logger, client)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client
}

type logger struct{}

func (l *logger) Printf(_ context.Context, _ string, _ ...any) {}

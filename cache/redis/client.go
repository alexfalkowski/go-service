package redis

import (
	"context"

	gz "github.com/alexfalkowski/go-service/cache/redis/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/cache/redis/telemetry/tracer"
	gr "github.com/alexfalkowski/go-service/redis"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ClientParams for redis.
type ClientParams struct {
	fx.In

	Lifecycle   fx.Lifecycle
	RingOptions *redis.RingOptions
	Tracer      trace.Tracer
	Logger      *zap.Logger
}

// NewClient for redis.
func NewClient(params ClientParams) gr.Client {
	r := params.RingOptions
	if r == nil {
		return gr.NewNoopClient()
	}

	redis.SetLogger(gz.NewLogger(params.Logger))

	var client gr.Client = redis.NewRing(r)
	client = gz.NewClient(params.Logger, client)
	client = tracer.NewClient(params.Tracer, client)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return client.Close()
		},
	})

	return client
}

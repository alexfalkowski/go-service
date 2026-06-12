package driver

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/cache/telemetry"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/redis/go-redis/v9"
	notifications "github.com/redis/go-redis/v9/maintnotifications"
)

func newRedisDriver(params DriverParams) (Driver, error) {
	data, err := params.FS.ReadSource(params.Config.Options["url"].(string))
	if err != nil {
		return nil, err
	}

	opts, err := redis.ParseURL(bytes.String(data))
	if err != nil {
		return nil, ErrInvalidURL
	}

	opts.MaintNotificationsConfig = &notifications.Config{
		Mode: notifications.ModeDisabled,
	}
	if params.Logger != nil {
		redis.SetLogger(redisLogger{logger: params.Logger})
	}

	redisClient := redis.NewClient(opts)
	if tracer.IsEnabled() {
		runtime.Must(telemetry.InstrumentTracing(redisClient))
	}

	var metricsClose chan struct{}
	if metrics.IsEnabled() {
		metricsClose = make(chan struct{})
		runtime.Must(telemetry.InstrumentMetrics(redisClient, metricsClose))
	}

	params.Lifecycle.Append(di.Hook{
		OnStop: func(context.Context) error {
			if metricsClose != nil {
				close(metricsClose)
			}

			return redisClient.Close()
		},
	})

	return &redisDriver{client: redisClient}, nil
}

type redisDriver struct {
	client *redis.Client
}

func (d *redisDriver) Delete(ctx context.Context, key string) error {
	return d.client.Del(ctx, key).Err()
}

func (d *redisDriver) Fetch(ctx context.Context, key string) (string, error) {
	value, err := d.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrMissing
	}

	return value, err
}

func (d *redisDriver) Flush(ctx context.Context) error {
	return d.client.FlushDB(ctx).Err()
}

func (d *redisDriver) Save(ctx context.Context, key, value string, lifetime time.Duration) error {
	return d.client.Set(ctx, key, value, lifetime.Duration()).Err()
}

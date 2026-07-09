// Package redis provides the internal Redis cache driver.
package redis

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/cache/config"
	drivererrors "github.com/alexfalkowski/go-service/v2/cache/driver/errors"
	"github.com/alexfalkowski/go-service/v2/cache/telemetry"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/redis/go-redis/v9"
	notifications "github.com/redis/go-redis/v9/maintnotifications"
)

// NewDriver constructs a Redis cache driver.
func NewDriver(lc di.Lifecycle, fs *os.FS, cfg *config.Config, log *logger.Logger) (*Driver, error) {
	data, err := fs.ReadSource(cfg.Options["url"].(string))
	if err != nil {
		return nil, err
	}

	opts, err := redis.ParseURL(bytes.String(data))
	if err != nil {
		return nil, drivererrors.ErrInvalidURL
	}

	opts.MaintNotificationsConfig = &notifications.Config{
		Mode: notifications.ModeDisabled,
	}
	if log != nil {
		redis.SetLogger(redisLogger{logger: log})
	}

	client := redis.NewClient(opts)
	if tracer.IsEnabled() {
		runtime.Must(telemetry.InstrumentTracing(client))
	}

	var metricsClose chan struct{}
	if metrics.IsEnabled() {
		metricsClose = make(chan struct{})
		runtime.Must(telemetry.InstrumentMetrics(client, metricsClose))
	}

	lc.Append(di.Hook{
		OnStop: func(context.Context) error {
			if metricsClose != nil {
				close(metricsClose)
			}

			return client.Close()
		},
	})

	return &Driver{client: client}, nil
}

// Driver is a Redis cache driver.
type Driver struct {
	client *redis.Client
}

// Delete removes the cached key.
func (d *Driver) Delete(ctx context.Context, key string) error {
	return d.client.Del(ctx, key).Err()
}

// Get retrieves the cached value for key.
func (d *Driver) Get(ctx context.Context, key string) (string, error) {
	value, err := d.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", drivererrors.ErrMissing
	}

	return value, err
}

// Flush clears the entire selected Redis database.
//
// This uses Redis FLUSHDB, so it removes keys that were not created through the
// go-service cache facade when they share the same selected database.
func (d *Driver) Flush(ctx context.Context) error {
	return d.client.FlushDB(ctx).Err()
}

// Ping verifies Redis connectivity.
func (d *Driver) Ping(ctx context.Context) error {
	return d.client.Ping(ctx).Err()
}

// Save stores value under key for the provided lifetime.
func (d *Driver) Save(ctx context.Context, key, value string, lifetime time.Duration) error {
	return d.client.Set(ctx, key, value, lifetime.Duration()).Err()
}

// GetOrSave atomically returns the value stored for key, or stores value when the key is absent.
//
// It issues a single SET ... NX GET command, so it requires Redis 7.0 or later.
func (d *Driver) GetOrSave(ctx context.Context, key, value string, lifetime time.Duration) (string, bool, error) {
	existing, err := d.client.SetArgs(ctx, key, value, redis.SetArgs{
		Mode: "NX",
		Get:  true,
		TTL:  lifetime.Duration(),
	}).Result()
	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	return existing, true, nil
}

package driver

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/redis/go-redis/v9"
)

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

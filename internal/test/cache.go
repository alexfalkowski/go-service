package test

import (
	"github.com/alexfalkowski/go-service/v2/cache"
	"github.com/alexfalkowski/go-service/v2/cache/cacher"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
)

type Cache struct {
	Value string
}

func (c *Cache) Contains(_ string) bool {
	return true
}

func (c *Cache) Delete(_ string) error {
	return nil
}

func (c *Cache) Fetch(_ string) (string, error) {
	return c.Value, nil
}

func (c *Cache) FetchMulti(_ []string) map[string]string {
	return map[string]string{}
}

func (c *Cache) Flush() error {
	return nil
}

func (c *Cache) Save(_, _ string, _ time.Duration) error {
	return nil
}

type ErrCache struct{}

func (*ErrCache) Contains(_ string) bool {
	return true
}

func (*ErrCache) Delete(_ string) error {
	return ErrFailed
}

func (*ErrCache) Fetch(_ string) (string, error) {
	return strings.Empty, ErrFailed
}

func (*ErrCache) FetchMulti(_ []string) map[string]string {
	return map[string]string{}
}

func (*ErrCache) Flush() error {
	return nil
}

func (*ErrCache) Save(_, _ string, _ time.Duration) error {
	return ErrFailed
}

func redisCache(lc di.Lifecycle, logger *logger.Logger, meter *metrics.Meter, tracer *tracer.Config) cacher.Cache {
	cfg := NewCacheConfig("redis", "snappy", "json", "redis")

	driver, err := driver.NewDriver(FS, cfg)
	runtime.Must(err)

	params := cache.CacheParams{
		Lifecycle:  lc,
		Config:     cfg,
		Compressor: Compressor,
		Encoder:    Encoder,
		Pool:       Pool,
		Driver:     driver,
		Tracer:     NewTracer(lc, tracer),
		Logger:     logger,
		Meter:      meter,
	}

	return cache.NewCache(params)
}

package cache

import (
	"bytes"
	"context"
	"encoding/base64"
	"time"

	"github.com/alexfalkowski/go-service/cache/config"
	cz "github.com/alexfalkowski/go-service/cache/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/cache/telemetry/metrics"
	"github.com/alexfalkowski/go-service/cache/telemetry/tracer"
	"github.com/alexfalkowski/go-service/compress"
	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/sync"
	"github.com/alexfalkowski/go-service/types/ptr"
	"github.com/faabiosr/cachego"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Get a value from key.
func Get[T any](ctx context.Context, cache config.Cache, key string) (*T, error) {
	value := ptr.Zero[T]()
	err := cache.Get(ctx, key, value)

	return value, err
}

// Persist a value to the key with a TTL.
func Persist[T any](ctx context.Context, cache config.Cache, key string, value *T, ttl time.Duration) error {
	return cache.Persist(ctx, key, value, ttl)
}

// Params for cache.
type Params struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Config     *config.Config
	Encoder    *encoding.Map
	Pool       *sync.BufferPool
	Compressor *compress.Map
	Cache      cachego.Cache
	Tracer     trace.Tracer
	Logger     *zap.Logger
	Meter      metric.Meter
}

// New from config.
func New(params Params) (config.Cache, error) {
	if !config.IsEnabled(params.Config) {
		return nil, nil
	}

	cmp := params.Compressor.Get(params.Config.Compressor)
	enc := params.Encoder.Get(params.Config.Encoder)

	var cache config.Cache = &Cache{cmp: cmp, enc: enc, pool: params.Pool, cache: params.Cache}
	cache = tracer.NewCache(params.Config.Kind, params.Tracer, cache)
	cache = cz.NewCache(params.Config.Kind, params.Logger, cache)
	cache = metrics.NewCache(params.Config.Kind, params.Meter, cache)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return cache.Close(ctx)
		},
	})

	return cache, nil
}

// Cache allows marshaling and compressing items to the cache.
type Cache struct {
	enc   encoding.Encoder
	pool  *sync.BufferPool
	cmp   compress.Compressor
	cache cachego.Cache
}

// Close the cache.
func (c *Cache) Close(_ context.Context) error {
	return c.cache.Flush()
}

// Remove a cached key.
func (c *Cache) Remove(_ context.Context, key string) error {
	return c.cache.Delete(key)
}

// Get a cached value.
func (c *Cache) Get(_ context.Context, key string, value any) error {
	val, err := c.cache.Fetch(key)
	if err != nil {
		return err
	}

	return c.decode(val, value)
}

// Persist a value with key and TTL.
func (c *Cache) Persist(_ context.Context, key string, value any, ttl time.Duration) error {
	enc, err := c.encode(value)
	if err != nil {
		return err
	}

	return c.cache.Save(key, enc, ttl)
}

func (c *Cache) encode(value any) (string, error) {
	buf := c.pool.Get()
	defer c.pool.Put(buf)

	if err := c.enc.Encode(buf, value); err != nil {
		return "", err
	}

	cmp := c.cmp.Compress(buf.Bytes())

	return base64.StdEncoding.EncodeToString(cmp), nil
}

func (c *Cache) decode(value string, field any) error {
	data, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return err
	}

	data, err = c.cmp.Decompress(data)
	if err != nil {
		return err
	}

	return c.enc.Decode(bytes.NewReader(data), field)
}

package cache

import (
	"context"

	"github.com/alexfalkowski/go-service/bytes"
	"github.com/alexfalkowski/go-service/cache/cacheable"
	"github.com/alexfalkowski/go-service/cache/config"
	"github.com/alexfalkowski/go-service/cache/driver"
	cl "github.com/alexfalkowski/go-service/cache/telemetry/logger"
	cm "github.com/alexfalkowski/go-service/cache/telemetry/metrics"
	ct "github.com/alexfalkowski/go-service/cache/telemetry/tracer"
	"github.com/alexfalkowski/go-service/compress"
	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/encoding/base64"
	"github.com/alexfalkowski/go-service/sync"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/time"
	"go.uber.org/fx"
)

// Params for cache.
type Params struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Config     *config.Config
	Encoder    *encoding.Map
	Pool       *sync.BufferPool
	Compressor *compress.Map
	Driver     driver.Driver
	Tracer     *tracer.Tracer
	Logger     *logger.Logger
	Meter      *metrics.Meter
}

// NewCache from config.
func NewCache(params Params) cacheable.Interface {
	if !config.IsEnabled(params.Config) {
		return nil
	}

	cmp := params.Compressor.Get(params.Config.Compressor)
	enc := params.Encoder.Get(params.Config.Encoder)

	var cache cacheable.Interface = &Cache{compressor: cmp, encoder: enc, pool: params.Pool, driver: params.Driver}

	if params.Tracer != nil {
		cache = ct.NewCache(params.Config.Kind, params.Tracer, cache)
	}

	if params.Logger != nil {
		cache = cl.NewCache(params.Config.Kind, params.Logger, cache)
	}

	if params.Meter != nil {
		cache = cm.NewCache(params.Config.Kind, params.Meter, cache)
	}

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return cache.Close(ctx)
		},
	})

	return cache
}

// Cache allows marshaling and compressing items to the cache.
type Cache struct {
	encoder    encoding.Encoder
	pool       *sync.BufferPool
	compressor compress.Compressor
	driver     driver.Driver
}

// Close the cache.
func (c *Cache) Close(_ context.Context) error {
	return c.driver.Flush()
}

// Remove a cached key.
func (c *Cache) Remove(_ context.Context, key string) (bool, error) {
	if !c.driver.Contains(key) {
		return false, nil
	}

	return true, c.driver.Delete(key)
}

// Get a cached value.
func (c *Cache) Get(_ context.Context, key string, value any) (bool, error) {
	if !c.driver.Contains(key) {
		return false, nil
	}

	val, err := c.driver.Fetch(key)
	if err != nil {
		return false, err
	}

	return true, c.decode(val, value)
}

// Persist a value with key and TTL.
func (c *Cache) Persist(_ context.Context, key string, value any, ttl time.Duration) error {
	enc, err := c.encode(value)
	if err != nil {
		return err
	}

	return c.driver.Save(key, enc, ttl)
}

func (c *Cache) encode(value any) (string, error) {
	var data []byte

	switch kind := value.(type) {
	case *[]byte:
		data = *kind
	case *bytes.Buffer:
		data = kind.Bytes()
	default:
		buf := c.pool.Get()
		defer c.pool.Put(buf)

		if err := c.encoder.Encode(buf, value); err != nil {
			return "", err
		}

		data = buf.Bytes()
	}

	compressed := c.compressor.Compress(data)
	encoded := base64.Encode(compressed)

	return encoded, nil
}

func (c *Cache) decode(value string, field any) error {
	decoded, err := base64.Decode(value)
	if err != nil {
		return err
	}

	decompressed, err := c.compressor.Decompress(decoded)
	if err != nil {
		return err
	}

	switch kind := field.(type) {
	case *[]byte:
		*kind = decompressed

		return nil
	case *bytes.Buffer:
		_, _ = kind.Write(decompressed)

		return nil
	default:
		return c.encoder.Decode(bytes.NewReader(decompressed), field)
	}
}

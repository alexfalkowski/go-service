package cache

import (
	"bytes"
	"context"
	"encoding/base64"

	"github.com/alexfalkowski/go-service/compress"
	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/sync"
	"github.com/faabiosr/cachego"
	"github.com/faabiosr/cachego/redis"
	cs "github.com/faabiosr/cachego/sync"
	client "github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

// Params for cache.
type Params struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Config     *Config
	Encoder    *encoding.Map
	Pool       *sync.BufferPool
	Compressor *compress.Map
}

// New from config.
func New(params Params) (cache *Cache, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Prefix("cache", runtime.ConvertRecover(r))
		}
	}()

	cache = &Cache{
		Compressor: params.Compressor.Get(params.Config.Compressor),
		Encoder:    params.Encoder.Get(params.Config.Encoder),
		Pool:       params.Pool,
	}

	switch params.Config.Kind {
	case "redis":
		url, err := os.ReadFile(params.Config.Options["url"].(string))
		runtime.Must(err)

		opts, err := client.ParseURL(url)
		runtime.Must(err)

		cache.Cache = redis.New(client.NewClient(opts))
	default:
		cache.Cache = cs.New()
	}

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return cache.Flush()
		},
	})

	return
}

// Cache allows marshaling and compressing items to the cache.
type Cache struct {
	Encoder    encoding.Encoder
	Pool       *sync.BufferPool
	Compressor compress.Compressor
	cachego.Cache
}

// CreateValue encodes, compresses and base64 the value.
func (c *Cache) EncodeValue(value any) (string, error) {
	buf := c.Pool.Get()
	defer c.Pool.Put(buf)

	if err := c.Encoder.Encode(buf, value); err != nil {
		return "", err
	}

	cmp := c.Compressor.Compress(buf.Bytes())

	return base64.StdEncoding.EncodeToString(cmp), nil
}

// DecodeValue base64, uncompresses and decodes.
func (c *Cache) DecodeValue(value string, field any) error {
	data, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return err
	}

	data, err = c.Compressor.Decompress(data)
	if err != nil {
		return err
	}

	return c.Encoder.Decode(bytes.NewReader(data), field)
}

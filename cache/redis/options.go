package redis

import (
	"bytes"
	"time"

	"github.com/alexfalkowski/go-service/compress"
	"github.com/alexfalkowski/go-service/encoding"
	gr "github.com/alexfalkowski/go-service/redis"
	"github.com/alexfalkowski/go-service/sync"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

// OptionsParams for redis.
type OptionsParams struct {
	fx.In

	Client     gr.Client
	Config     *Config
	Encoder    *encoding.Map
	Pool       *sync.BufferPool
	Compressor *compress.Map
}

// NewOptions for redis.
func NewOptions(params OptionsParams) (*cache.Options, error) {
	if !IsEnabled(params.Config) {
		opts := &cache.Options{
			StatsEnabled: true,
			LocalCache:   cache.NewTinyLFU(1, time.Minute),
		}

		return opts, nil
	}

	enc := params.Encoder.Get(params.Config.Encoder)
	cmp := params.Compressor.Get(params.Config.Compressor)
	opts := &cache.Options{
		Redis:        params.Client,
		StatsEnabled: true,
		Marshal: func(v any) ([]byte, error) {
			buf := params.Pool.Get()
			defer params.Pool.Put(buf)

			if err := enc.Encode(buf, v); err != nil {
				return nil, err
			}

			return cmp.Compress(buf.Bytes()), nil
		},
		Unmarshal: func(b []byte, v any) error {
			d, err := cmp.Decompress(b)
			if err != nil {
				return err
			}

			return enc.Decode(bytes.NewReader(d), v)
		},
	}

	return opts, nil
}

// NewRingOptions for redis.
func NewRingOptions(cfg *Config) (*redis.RingOptions, error) {
	if cfg == nil {
		return nil, nil
	}

	u, err := cfg.GetURL()
	if err != nil {
		return nil, err
	}

	url, err := redis.ParseURL(u)
	if err != nil {
		return nil, err
	}

	opts := &redis.RingOptions{
		Addrs: cfg.Addresses,
		NewClient: func(*redis.Options) *redis.Client {
			return redis.NewClient(url)
		},
	}

	return opts, nil
}

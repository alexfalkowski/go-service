package redis

import (
	"github.com/alexfalkowski/go-service/cache/compressor"
	"github.com/alexfalkowski/go-service/cache/marshaller"
	gr "github.com/alexfalkowski/go-service/redis"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

// OptionsParams for redis.
type OptionsParams struct {
	fx.In

	Client     gr.Client
	Marshaller marshaller.Marshaller
	Compressor compressor.Compressor
}

// NewOptions for redis.
func NewOptions(params OptionsParams) *cache.Options {
	opts := &cache.Options{
		Redis:        params.Client,
		StatsEnabled: true,
		Marshal: func(v any) ([]byte, error) {
			d, err := params.Marshaller.Marshal(v)
			if err != nil {
				return nil, err
			}

			return params.Compressor.Compress(d), nil
		},
		Unmarshal: func(b []byte, v any) error {
			d, err := params.Compressor.Decompress(b)
			if err != nil {
				return err
			}

			return params.Marshaller.Unmarshal(d, v)
		},
	}

	return opts
}

// NewRingOptions for redis.
func NewRingOptions(cfg *Config) *redis.RingOptions {
	if cfg == nil {
		return nil
	}

	return &redis.RingOptions{
		Addrs:    cfg.Addresses,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	}
}

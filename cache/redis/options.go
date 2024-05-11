package redis

import (
	"time"

	"github.com/alexfalkowski/go-service/compressor"
	"github.com/alexfalkowski/go-service/marshaller"
	gr "github.com/alexfalkowski/go-service/redis"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

// OptionsParams for redis.
type OptionsParams struct {
	fx.In

	Client     gr.Client
	Config     *Config
	Marshaller *marshaller.Map
	Compressor *compressor.Map
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

	fm := params.Marshaller.Get(params.Config.Marshaller)
	cm := params.Compressor.Get(params.Config.Compressor)
	opts := &cache.Options{
		Redis:        params.Client,
		StatsEnabled: true,
		Marshal: func(v any) ([]byte, error) {
			d, err := fm.Marshal(v)
			if err != nil {
				return nil, err
			}

			return cm.Compress(d), nil
		},
		Unmarshal: func(b []byte, v any) error {
			d, err := cm.Decompress(b)
			if err != nil {
				return err
			}

			return fm.Unmarshal(d, v)
		},
	}

	return opts, nil
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

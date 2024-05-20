package redis

import (
	"time"

	"github.com/alexfalkowski/go-service/compressor"
	"github.com/alexfalkowski/go-service/marshaller"
	gr "github.com/alexfalkowski/go-service/redis"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
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
//
//nolint:nilnil
func NewRingOptions(cfg *Config) (*redis.RingOptions, error) {
	if cfg == nil {
		return nil, nil
	}

	u, err := cfg.GetURL()
	if err != nil {
		return nil, err
	}

	pu, err := redis.ParseURL(u)
	if err != nil {
		return nil, err
	}

	opts := &redis.RingOptions{
		Addrs: cfg.Addresses,
		NewClient: func(*redis.Options) *redis.Client {
			return redis.NewClient(pu)
		},
	}

	return opts, nil
}

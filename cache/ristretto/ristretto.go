package ristretto

import (
	"context"

	"github.com/alexfalkowski/go-service/cache/ristretto/metrics/prometheus"
	"github.com/alexfalkowski/go-service/version"
	"github.com/dgraph-io/ristretto"
	"go.uber.org/fx"
)

// CacheParams for ristretto.
type CacheParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *Config
	Version   version.Version
}

// NewCache for ristretto.
func NewCache(params CacheParams) (*ristretto.Cache, error) {
	rcfg := &ristretto.Config{
		NumCounters: params.Config.NumCounters,
		MaxCost:     params.Config.MaxCost,
		BufferItems: params.Config.BufferItems,
		Metrics:     true,
	}

	cache, err := ristretto.NewCache(rcfg)
	if err != nil {
		return nil, err
	}

	prometheus.Register(params.Lifecycle, cache, params.Version)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			cache.Close()

			return nil
		},
	})

	return cache, nil
}

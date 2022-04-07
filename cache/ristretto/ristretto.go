package ristretto

import (
	"context"

	"github.com/alexfalkowski/go-service/cache/ristretto/metrics/prometheus"
	"github.com/alexfalkowski/go-service/os"
	"github.com/dgraph-io/ristretto"
	"go.uber.org/fx"
)

// NewCache for ristretto.
func NewCache(lc fx.Lifecycle, cfg *Config) (*ristretto.Cache, error) {
	rcfg := &ristretto.Config{
		NumCounters: cfg.NumCounters,
		MaxCost:     cfg.MaxCost,
		BufferItems: cfg.BufferItems,
		Metrics:     true,
	}

	cache, err := ristretto.NewCache(rcfg)
	if err != nil {
		return nil, err
	}

	prometheus.Register(lc, os.ExecutableName(), cache)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			cache.Close()

			return nil
		},
	})

	return cache, nil
}

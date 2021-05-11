package ristretto

import (
	"context"

	"github.com/alexfalkowski/go-service/pkg/cache/ristretto/metrics/prometheus"
	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/dgraph-io/ristretto"
	"go.uber.org/fx"
)

// NewCache for ristretto.
func NewCache(lc fx.Lifecycle, cfg *config.Config, rcfg *ristretto.Config) (*ristretto.Cache, error) {
	cache, err := ristretto.NewCache(rcfg)
	if err != nil {
		return nil, err
	}

	if err := prometheus.Register(cfg, cache); err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			cache.Close()

			return nil
		},
	})

	return cache, nil
}

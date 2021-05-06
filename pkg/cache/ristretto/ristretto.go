package ristretto

import (
	"context"

	"github.com/dgraph-io/ristretto"
	"go.uber.org/fx"
)

// NewCache for ristretto.
func NewCache(lc fx.Lifecycle, cfg *ristretto.Config) (*ristretto.Cache, error) {
	cache, err := ristretto.NewCache(cfg)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			cache.Close()

			return nil
		},
	})

	return cache, err
}

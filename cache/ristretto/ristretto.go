package ristretto

import (
	"context"

	"github.com/alexfalkowski/go-service/errors"
	sr "github.com/alexfalkowski/go-service/ristretto"
	"github.com/dgraph-io/ristretto"
	"go.uber.org/fx"
)

// CacheParams for ristretto.
type CacheParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *Config
}

// NewCache for ristretto.
func NewCache(params CacheParams) (sr.Cache, error) {
	c := params.Config
	if c == nil {
		return sr.NewNoopCache(), nil
	}

	cfg := &ristretto.Config{
		NumCounters: c.NumCounters,
		MaxCost:     c.MaxCost,
		BufferItems: c.BufferItems,
		Metrics:     true,
	}

	ca, err := ristretto.NewCache(cfg)
	if err != nil {
		return nil, errors.Prefix("new ristretto options", err)
	}

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			ca.Close()

			return nil
		},
	})

	return &cache{Cache: ca}, nil
}

type cache struct {
	*ristretto.Cache
}

func (c *cache) Hits() uint64 {
	return c.Cache.Metrics.Hits()
}

func (c *cache) Misses() uint64 {
	return c.Cache.Metrics.Misses()
}

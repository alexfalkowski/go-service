package prometheus

import (
	"context"

	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/go-redis/cache/v8"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

// Register for prometheus.
func Register(lc fx.Lifecycle, cfg *config.Config, cache *cache.Cache) {
	collector := NewStatsCollector(cfg.AppName, cache)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return prometheus.Register(collector)
		},
		OnStop: func(ctx context.Context) error {
			prometheus.Unregister(collector)

			return nil
		},
	})
}

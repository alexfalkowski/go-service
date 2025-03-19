package transport

import (
	"context"

	"github.com/alexfalkowski/go-service/server"
	"go.uber.org/fx"
)

// Register all the transports.
func Register(lc fx.Lifecycle, servers []*server.Service) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			for _, s := range servers {
				s.Start()
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			for _, s := range servers {
				s.Stop(ctx)
			}

			return nil
		},
	})
}

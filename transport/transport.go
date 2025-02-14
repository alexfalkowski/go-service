package transport

import (
	"context"

	"go.uber.org/fx"
)

// Register all the transports.
func Register(lc fx.Lifecycle, servers []Server) {
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

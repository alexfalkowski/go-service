package transport

import (
	"context"

	"go.uber.org/fx"
)

// RegisterParams for transport.
type RegisterParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Servers   []Server
}

// Register all the transports.
func Register(params RegisterParams) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			for _, s := range params.Servers {
				s.Start()
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			for _, s := range params.Servers {
				s.Stop(ctx)
			}

			return nil
		},
	})
}

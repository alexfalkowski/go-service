package transport

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net/server"
)

// Register all the transports.
func Register(lc di.Lifecycle, servers []*server.Service) {
	lc.Append(di.Hook{
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

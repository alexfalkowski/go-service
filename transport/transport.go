package transport

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net/server"
)

// Register wires transport server services into the application lifecycle.
//
// It appends an Fx lifecycle hook that:
//
//   - OnStart: starts each provided `*server.Service` by calling `Start()`.
//   - OnStop: stops each provided `*server.Service` by calling `Stop(ctx)`.
//
// The `servers` slice is typically produced by `NewServers` and may contain services
// that are nil/disabled upstream. Callers should ensure the slice does not contain
// nil entries (for example by using helpers that append non-nil services).
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

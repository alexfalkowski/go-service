package server

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/di"
)

// Register wires server services into the application lifecycle.
//
// It appends an Fx lifecycle hook that:
//   - OnStart: starts each provided `*Service` by calling `Start()`.
//   - OnStop: stops each provided `*Service` by calling `Stop(ctx)`.
//
// Callers should ensure the slice does not contain nil entries.
func Register(lc di.Lifecycle, services []*Service) {
	lc.Append(di.Hook{
		OnStart: func(context.Context) error {
			for _, s := range services {
				s.Start()
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			for _, s := range services {
				s.Stop(ctx)
			}

			return nil
		},
	})
}

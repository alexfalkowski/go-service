package server

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/slices"
	"github.com/alexfalkowski/go-sync"
)

// RegisterParams defines dependencies for registering server services.
type RegisterParams struct {
	di.In

	// Lifecycle receives the server start and stop hooks.
	Lifecycle di.Lifecycle

	// Drain tracks whether shutdown has started.
	Drain *Drain

	// Services are the server services managed by the lifecycle.
	Services []*Service
}

// Register wires server services into the application lifecycle.
//
// It appends an Fx lifecycle hook that:
//   - OnStart: starts each provided *[Service] by calling `Start()`.
//   - OnStop: marks the drain state, stops each provided *[Service] by calling `Stop(ctx)`,
//     and returns the joined shutdown errors.
//
// Callers should ensure the slice does not contain nil entries.
func Register(params RegisterParams) {
	services := slices.Clone(params.Services)

	params.Lifecycle.Append(di.Hook{
		OnStart: func(context.Context) error {
			for _, s := range services {
				s.Start()
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			params.Drain.Start()

			group := sync.ErrorsGroup{}
			group.SetLimit(len(services))

			for _, s := range services {
				group.Go(func() error {
					return s.Stop(ctx)
				})
			}

			return group.Wait()
		},
	})
}

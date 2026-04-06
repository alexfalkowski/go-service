package server

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/errors"
)

// Register wires server services into the application lifecycle.
//
// It appends an Fx lifecycle hook that:
//   - OnStart: starts each provided `*Service` by calling `Start()`.
//   - OnStop: stops each provided `*Service` by calling `Stop(ctx)` and returns the joined shutdown errors.
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
			errs := make([]error, 0, len(services))
			for _, s := range services {
				errs = append(errs, s.Stop(ctx))
			}

			return errors.Join(errs...)
		},
	})
}

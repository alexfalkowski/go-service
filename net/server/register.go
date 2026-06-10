package server

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/slices"
	"github.com/alexfalkowski/go-sync"
)

// Register wires server services into the application lifecycle.
//
// It appends an Fx lifecycle hook that:
//   - OnStart: starts each provided *[Service] by calling `Start()`.
//   - OnStop: stops each provided *[Service] by calling `Stop(ctx)` and returns the joined shutdown errors.
//
// Callers should ensure the slice does not contain nil entries.
func Register(lc di.Lifecycle, services []*Service) {
	services = slices.Clone(services)

	lc.Append(di.Hook{
		OnStart: func(context.Context) error {
			for _, s := range services {
				s.Start()
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			errs := make([]error, len(services))
			var wg sync.WaitGroup
			for i, s := range services {
				wg.Go(func() {
					errs[i] = s.Stop(ctx)
				})
			}
			wg.Wait()

			return errors.Join(errs...)
		},
	})
}

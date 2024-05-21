package cmd

import (
	"context"

	"github.com/alexfalkowski/go-service/runtime"
	"go.uber.org/fx"
)

// StartFn for cmd.
type StartFn func(ctx context.Context)

// Start for cmd.
func Start(lc fx.Lifecycle, fn StartFn) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = runtime.ConvertRecover(r)
				}
			}()

			fn(ctx)

			return
		},
	})
}

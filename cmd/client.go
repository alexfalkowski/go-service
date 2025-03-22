package cmd

import (
	"context"

	"go.uber.org/fx"
)

// Start for cmd.
func Start(lc fx.Lifecycle, fn func(ctx context.Context) error) {
	lc.Append(fx.Hook{
		OnStart: fn,
	})
}

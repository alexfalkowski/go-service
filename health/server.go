package health

import (
	"context"

	"github.com/alexfalkowski/go-health/server"
	"go.uber.org/fx"
)

// NewServer for health.
func NewServer(lc fx.Lifecycle, regs Registrations) *server.Server {
	server := server.NewServer()
	server.Register(regs...)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			server.Start()

			return nil
		},
		OnStop: func(_ context.Context) error {
			server.Stop()

			return nil
		},
	})

	return server
}

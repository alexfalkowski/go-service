package health

import (
	"context"

	"github.com/alexfalkowski/go-health/server"
	"go.uber.org/fx"
)

// NewServer for health.
func NewServer(lc fx.Lifecycle, regs Registrations) *server.Server {
	s := server.NewServer()

	s.Register(regs...)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			s.Start()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.Stop()

			return nil
		},
	})

	return s
}

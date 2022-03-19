package health

import (
	"context"

	"github.com/alexfalkowski/go-health/server"
	"go.uber.org/fx"
)

// NewServer for health.
func NewServer(lc fx.Lifecycle, regs Registrations) (*server.Server, error) {
	s := server.NewServer()

	if err := s.Register(regs...); err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return s.Start()
		},
		OnStop: func(ctx context.Context) error {
			return s.Stop()
		},
	})

	return s, nil
}

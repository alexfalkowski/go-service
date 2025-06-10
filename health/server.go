package health

import (
	"context"

	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-service/v2/di"
)

// NewServer for health.
func NewServer(lc di.Lifecycle, regs Registrations) *server.Server {
	server := server.NewServer()
	server.Register(regs...)

	lc.Append(di.Hook{
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

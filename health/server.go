package health

import (
	health "github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
)

// NewServer constructs a go-health server and wires it into the application lifecycle.
//
// The server is started on lifecycle start and stopped on lifecycle stop.
func NewServer(lc di.Lifecycle) *health.Server {
	server := health.NewServer()

	lc.Append(di.Hook{
		OnStart: func(_ context.Context) error {
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

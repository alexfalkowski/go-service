package health

import (
	health "github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
)

// NewServer constructs a go-health server and wires it into the application lifecycle.
//
// Lifecycle behavior:
//   - OnStart: starts the underlying go-health server.
//   - OnStop: stops the underlying go-health server.
//
// This makes the health server's lifetime match the Fx application lifetime, without requiring
// callers to manually start/stop it. The returned server can then be configured by other wiring
// (for example by registering health check registrations/endpoints provided by go-health).
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

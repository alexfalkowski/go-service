package health

import (
	health "github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-service/v2/di"
)

// NewServer constructs a go-health server and wires it into the application lifecycle.
//
// It starts the underlying go-health server on application start and stops it on application stop.
// This makes the health server's lifetime match the Fx application lifetime, without requiring callers
// to manually start/stop it. The returned server is the shared observer/registration store used by
// transport-specific health endpoint wiring.
func NewServer(lc di.Lifecycle) *health.Server {
	server := health.NewServer()

	lc.Append(di.Hook{
		OnStart: server.Start,
		OnStop:  server.Stop,
	})

	return server
}

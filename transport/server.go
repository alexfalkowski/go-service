package transport

import (
	"context"

	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
	"go.uber.org/fx"
)

// ServersParams for transport.
type ServersParams struct {
	fx.In

	HTTP  *http.Server
	GRPC  *grpc.Server
	Debug *debug.Server
}

// NewServers for transport.
func NewServers(params ServersParams) []Server {
	servers := []Server{}

	if params.HTTP != nil {
		servers = append(servers, params.HTTP)
	}

	if params.GRPC != nil {
		servers = append(servers, params.GRPC)
	}

	if params.Debug != nil {
		servers = append(servers, params.Debug)
	}

	return servers
}

// Server for transport.
type Server interface {
	// Start a server.
	Start()

	// Stop a server.
	Stop(ctx context.Context)
}

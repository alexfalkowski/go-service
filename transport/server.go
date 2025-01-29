package transport

import (
	"context"

	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/proxy"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
	"go.uber.org/fx"
)

type (
	// Server for transport.
	Server interface {
		Start()
		Stop(ctx context.Context)
	}

	// ServersParams for transport.
	ServersParams struct {
		fx.In
		HTTP  *http.Server
		GRPC  *grpc.Server
		DEBUG *debug.Server
		PROXY *proxy.Server
	}
)

// NewServers for transport.
func NewServers(params ServersParams) []Server {
	return []Server{params.HTTP, params.GRPC, params.DEBUG, params.PROXY}
}

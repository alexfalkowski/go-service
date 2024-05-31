package transport

import (
	"context"

	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/transport/ssh"
	"go.uber.org/fx"
)

// Server for transport.
type Server interface {
	Start()
	Stop(ctx context.Context)
}

// ServerParams for transport.
type ServerParams struct {
	fx.In

	Debug *debug.Server
	GRPC  *grpc.Server
	HTTP  *http.Server
	SSH   *ssh.Server
}

// NewServers for transport.
func NewServers(params ServerParams) []Server {
	return []Server{params.Debug, params.GRPC, params.HTTP, params.SSH}
}

package transport

import (
	"context"

	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
)

// Server for transport.
type Server interface {
	Start() error
	Stop(ctx context.Context) error
}

// NewServers for transport.
func NewServers(http *http.Server, grpc *grpc.Server, debug *debug.Server) []Server {
	return []Server{http, grpc, debug}
}

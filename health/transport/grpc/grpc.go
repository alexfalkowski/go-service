package grpc

import (
	"github.com/alexfalkowski/go-service/transport/grpc"
	"go.uber.org/fx"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

// RegisterParams for gRPC.
type RegisterParams struct {
	fx.In

	GRPC   *grpc.Server
	Server *Server
}

// Register health for gRPC.
func Register(params RegisterParams) {
	if params.GRPC == nil || params.Server == nil {
		return
	}

	health.RegisterHealthServer(params.GRPC.ServiceRegistrar(), params.Server)
}

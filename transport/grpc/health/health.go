package health

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

// RegisterParams for gRPC.
type RegisterParams struct {
	di.In
	Registrar grpc.ServiceRegistrar
	Server    *Server
}

// Register health for gRPC.
func Register(params RegisterParams) {
	if params.Registrar == nil || params.Server == nil {
		return
	}

	health.RegisterHealthServer(params.Registrar, params.Server)
}

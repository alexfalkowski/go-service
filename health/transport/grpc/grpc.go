package grpc

import (
	"github.com/alexfalkowski/go-service/transport/grpc"
	"go.uber.org/fx"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

// RegisterParams for gRPC.
type RegisterParams struct {
	fx.In

	Server   *grpc.Server
	Observer *Observer `optional:"true"`
}

// Register health for gRPC.
func Register(params RegisterParams) {
	if params.Server == nil || params.Observer == nil {
		return
	}

	health.RegisterHealthServer(params.Server.Server(), &server{ob: params.Observer})
}

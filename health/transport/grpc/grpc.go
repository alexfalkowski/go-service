package grpc

import (
	"github.com/alexfalkowski/go-service/transport/grpc"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

// Register health for gRPC.
func Register(srv *grpc.Server, ob *Observer) {
	health.RegisterHealthServer(srv.Server, &server{ob: ob})
}

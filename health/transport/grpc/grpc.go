package grpc

import (
	"google.golang.org/grpc"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

// Register health for gRPC.
func Register(srv *grpc.Server, ob *Observer) {
	health.RegisterHealthServer(srv, &server{ob: ob})
}

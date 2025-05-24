package grpc

import (
	"github.com/alexfalkowski/go-service/v2/time"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// NewServer for grpc.
func NewServer(timeout time.Duration, opts ...grpc.ServerOption) *grpc.Server {
	options := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             timeout,
			PermitWithoutStream: true,
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     timeout,
			MaxConnectionAge:      timeout,
			MaxConnectionAgeGrace: timeout,
			Time:                  timeout,
			Timeout:               timeout,
		}),
	}
	options = append(options, opts...)

	server := grpc.NewServer(options...)
	reflection.Register(server)

	return server
}

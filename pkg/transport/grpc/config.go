package grpc

import "github.com/alexfalkowski/go-service/pkg/transport/grpc/retry"

// Config for gRPC.
type Config struct {
	Port  string       `yaml:"port"`
	Retry retry.Config `yaml:"retry"`
}

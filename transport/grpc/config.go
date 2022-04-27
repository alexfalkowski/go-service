package grpc

import (
	"github.com/alexfalkowski/go-service/transport/grpc/retry"
)

// Config for gRPC.
type Config struct {
	Port      string       `yaml:"port"`
	Retry     retry.Config `yaml:"retry"`
	UserAgent string       `yaml:"user_agent"`
}

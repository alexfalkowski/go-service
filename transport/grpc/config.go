package grpc

import (
	"github.com/alexfalkowski/go-service/transport/grpc/retry"
)

// Config for gRPC.
type Config struct {
	Port      string       `yaml:"port"`
	UserAgent string       `yaml:"user_agent"`
	Retry     retry.Config `yaml:"retry"`
}

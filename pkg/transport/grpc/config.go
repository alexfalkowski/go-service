package grpc

import (
	"github.com/alexfalkowski/go-service/pkg/transport/grpc/ratelimit"
	"github.com/alexfalkowski/go-service/pkg/transport/grpc/retry"
)

// Config for gRPC.
type Config struct {
	Port      string           `yaml:"port"`
	UserAgent string           `yaml:"user_agent"`
	Retry     retry.Config     `yaml:"retry"`
	RateLimit ratelimit.Config `yaml:"rate_limit"`
}

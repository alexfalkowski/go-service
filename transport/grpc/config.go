package grpc

import (
	"github.com/alexfalkowski/go-service/transport/grpc/retry"
)

// Config for gRPC.
type Config struct {
	Retry     retry.Config `yaml:"retry" json:"retry" toml:"retry"`
	UserAgent string       `yaml:"user_agent" json:"user_agent" toml:"user_agent"`
}

package grpc

import (
	"github.com/alexfalkowski/go-service/server"
)

// Config for gRPC.
type Config struct {
	server.Config `yaml:",inline" json:",inline" toml:",inline"`
}

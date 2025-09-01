package grpc

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/limiter"
)

// NewServerLimiter for gRPC.
func NewServerLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *Config) (*limiter.Server, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	return limiter.NewServerLimiter(lc, keys, cfg.Limiter)
}

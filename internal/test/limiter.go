package test

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/limiter"
	gl "github.com/alexfalkowski/go-service/v2/transport/grpc/limiter"
	hl "github.com/alexfalkowski/go-service/v2/transport/http/limiter"
)

// LimiterKeyMap is the shared limiter key registry used by client and server limiter helpers.
var LimiterKeyMap = limiter.NewKeyMap()

// NewHTTPClientLimiter returns an HTTP client limiter and any construction error.
func NewHTTPClientLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) (*hl.Client, error) {
	return newHTTPClientLimiter(lc, keys, cfg)
}

// NewHTTPServerLimiter returns an HTTP server limiter and any construction error.
func NewHTTPServerLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) (*hl.Server, error) {
	return newHTTPServerLimiter(lc, keys, cfg)
}

// NewGRPCClientLimiter returns a gRPC client limiter and any construction error.
func NewGRPCClientLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) (*gl.Client, error) {
	return newGRPCClientLimiter(lc, keys, cfg)
}

// NewGRPCServerLimiter returns a gRPC server limiter and any construction error.
func NewGRPCServerLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) (*gl.Server, error) {
	return newGRPCServerLimiter(lc, keys, cfg)
}

func newHTTPClientLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) (*hl.Client, error) {
	return hl.NewClientLimiter(lc, keys, cfg)
}

func newHTTPServerLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) (*hl.Server, error) {
	return hl.NewServerLimiter(lc, keys, cfg)
}

func newGRPCClientLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) (*gl.Client, error) {
	return gl.NewClientLimiter(lc, keys, cfg)
}

func newGRPCServerLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) (*gl.Server, error) {
	return gl.NewServerLimiter(lc, keys, cfg)
}

package test

import (
	"github.com/alexfalkowski/go-service/v2/di"
	grpc "github.com/alexfalkowski/go-service/v2/transport/grpc/limiter"
	http "github.com/alexfalkowski/go-service/v2/transport/http/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/limiter"
)

// LimiterKeyMap is the shared limiter key registry used by client and server limiter helpers.
var LimiterKeyMap = limiter.NewKeyMap()

// NewHTTPClientLimiter returns an HTTP client limiter and any construction error.
func NewHTTPClientLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) (*http.Client, error) {
	return http.NewClientLimiter(lc, keys, cfg)
}

// NewHTTPServerLimiter returns an HTTP server limiter and any construction error.
func NewHTTPServerLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) (*http.Server, error) {
	return http.NewServerLimiter(lc, keys, cfg)
}

// NewGRPCClientLimiter returns a gRPC client limiter and any construction error.
func NewGRPCClientLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) (*grpc.Client, error) {
	return grpc.NewClientLimiter(lc, keys, cfg)
}

// NewGRPCServerLimiter returns a gRPC server limiter and any construction error.
func NewGRPCServerLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) (*grpc.Server, error) {
	return grpc.NewServerLimiter(lc, keys, cfg)
}

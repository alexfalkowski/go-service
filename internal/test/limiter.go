package test

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/runtime"
	gl "github.com/alexfalkowski/go-service/v2/transport/grpc/limiter"
	hl "github.com/alexfalkowski/go-service/v2/transport/http/limiter"
)

// LimiterKeyMap is the shared limiter key registry used by client and server limiter helpers.
var LimiterKeyMap = limiter.NewKeyMap()

// NewHTTPClientLimiter returns an HTTP client limiter and panics on construction errors.
func NewHTTPClientLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) *hl.Client {
	c, err := newHTTPClientLimiter(lc, keys, cfg)
	runtime.Must(err)

	return c
}

// NewHTTPServerLimiter returns an HTTP server limiter and panics on construction errors.
func NewHTTPServerLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) *hl.Server {
	s, err := newHTTPServerLimiter(lc, keys, cfg)
	runtime.Must(err)

	return s
}

// NewGRPCClientLimiter returns a gRPC client limiter and panics on construction errors.
func NewGRPCClientLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) *gl.Client {
	c, err := newGRPCClientLimiter(lc, keys, cfg)
	runtime.Must(err)

	return c
}

// NewGRPCServerLimiter returns a gRPC server limiter and panics on construction errors.
func NewGRPCServerLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) *gl.Server {
	s, err := newGRPCServerLimiter(lc, keys, cfg)
	runtime.Must(err)

	return s
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

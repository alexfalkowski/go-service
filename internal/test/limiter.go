package test

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/runtime"
	gl "github.com/alexfalkowski/go-service/v2/transport/grpc/limiter"
	hl "github.com/alexfalkowski/go-service/v2/transport/http/limiter"
)

// LimiterKeyMap for test.
var LimiterKeyMap = limiter.NewKeyMap()

// NewHTTPClientLimiter for test.
func NewHTTPClientLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) *hl.Client {
	c, err := hl.NewClientLimiter(lc, keys, cfg)
	runtime.Must(err)

	return c
}

// NewHTTPServerLimiter for test.
func NewHTTPServerLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) *hl.Server {
	s, err := hl.NewServerLimiter(lc, keys, cfg)
	runtime.Must(err)

	return s
}

// NewGRPCClientLimiter for test.
func NewGRPCClientLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) *gl.Client {
	c, err := gl.NewClientLimiter(lc, keys, cfg)
	runtime.Must(err)

	return c
}

// NewGRPCServerLimiter for test.
func NewGRPCServerLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *limiter.Config) *gl.Server {
	s, err := gl.NewServerLimiter(lc, keys, cfg)
	runtime.Must(err)

	return s
}

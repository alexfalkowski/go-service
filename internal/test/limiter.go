package test

import (
	"github.com/alexfalkowski/go-service/v2/context"
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

// AllowAllLimiter is a [github.com/alexfalkowski/go-service/v2/net/http/meta.RateLimiter] test double
// whose Take always allows.
type AllowAllLimiter struct{}

// Take always allows.
func (*AllowAllLimiter) Take(context.Context) (bool, string, error) {
	return true, "", nil
}

// CountingLimiter is a [github.com/alexfalkowski/go-service/v2/net/http/meta.RateLimiter] test double
// that allows exactly Remaining Take calls before denying every call after, used to verify per-message
// limiter charging deterministically without depending on transport/limiter's concrete type.
type CountingLimiter struct {
	Remaining int
}

// Take reports true and decrements Remaining while Remaining > 0, then reports false forever after.
func (l *CountingLimiter) Take(context.Context) (bool, string, error) {
	if l.Remaining <= 0 {
		return false, "", nil
	}

	l.Remaining--

	return true, "", nil
}

// ErrorLimiter is a [github.com/alexfalkowski/go-service/v2/net/http/meta.RateLimiter] test double
// whose Take always fails with ErrFailed.
type ErrorLimiter struct{}

// Take always returns ErrFailed.
func (ErrorLimiter) Take(context.Context) (bool, string, error) {
	return false, "", ErrFailed
}

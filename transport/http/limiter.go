package http

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/transport/http/limiter"
)

// NewServerLimiter for HTTP.
func NewServerLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *Config) (*limiter.Server, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}
	return limiter.NewServerLimiter(lc, keys, cfg.Limiter)
}

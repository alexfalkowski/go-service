package http

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/transport/http/limiter"
)

// NewServerLimiter constructs an HTTP server-side rate limiter when enabled.
//
// This is a small wiring helper that adapts the HTTP transport `Config` to the limiter package constructor.
// It is intended to be used by Fx module wiring.
//
// If cfg is disabled, it returns (nil, nil) so downstream wiring can treat rate limiting as not configured.
//
// The returned limiter is registered with the provided lifecycle. The `keys` map controls how request contexts
// are turned into limiter keys (for example, per user-agent or per client IP).
func NewServerLimiter(lc di.Lifecycle, keys limiter.KeyMap, cfg *Config) (*limiter.Server, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}
	return limiter.NewServerLimiter(lc, keys, cfg.Limiter)
}

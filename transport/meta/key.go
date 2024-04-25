package meta

import "github.com/alexfalkowski/go-service/limiter"

// NewKey for meta.
func NewKey(cfg *limiter.Config) limiter.KeyFunc {
	if !limiter.IsEnabled(cfg) {
		return limiter.NoKey
	}

	if cfg.Kind == "user-agent" {
		return UserAgent
	}

	return limiter.NoKey
}

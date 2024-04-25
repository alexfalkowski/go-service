package meta

import "github.com/alexfalkowski/go-service/limiter"

var keys = map[string]limiter.KeyFunc{
	"user-agent": UserAgent,
}

// NewKey for meta.
func NewKey(cfg *limiter.Config) limiter.KeyFunc {
	if !limiter.IsEnabled(cfg) {
		return nil
	}

	if u, ok := keys[cfg.Kind]; ok {
		return u
	}

	return nil
}

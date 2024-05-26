package limiter

import (
	"context"
	"errors"

	ge "github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

var (
	// ErrMissingKey for limiter.
	ErrMissingKey = errors.New("missing key")

	keys = map[string]KeyFunc{}
)

// RegisterKey with name and fn. Last register wins.
func RegisterKey(name string, fn KeyFunc) {
	keys[name] = fn
}

// New limiter.
func New(cfg *Config) (*limiter.Limiter, KeyFunc, error) {
	if !IsEnabled(cfg) {
		return nil, nil, nil
	}

	k, ok := keys[cfg.Kind]
	if !ok {
		return nil, nil, ErrMissingKey
	}

	rate, err := limiter.NewRateFromFormatted(cfg.Pattern)
	if err != nil {
		return nil, nil, ge.Prefix("new limiter", err)
	}

	return limiter.New(memory.NewStore(), rate), k, nil
}

// KeyFunc for the limiter.
type KeyFunc func(context.Context) meta.Valuer

package limiter

import (
	"context"
	"errors"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/time"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/memorystore"
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
func New(cfg *Config) (limiter.Store, KeyFunc, error) {
	if !IsEnabled(cfg) {
		return nil, nil, nil
	}

	k, ok := keys[cfg.Kind]
	if !ok {
		return nil, nil, ErrMissingKey
	}

	store, err := memorystore.New(&memorystore.Config{Tokens: cfg.Tokens, Interval: time.MustParseDuration(cfg.Interval)})
	if err != nil {
		return nil, nil, err
	}

	return store, k, nil
}

// KeyFunc for the limiter.
type KeyFunc func(context.Context) meta.Valuer

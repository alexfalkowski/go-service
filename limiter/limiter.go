package limiter

import (
	"context"
	"strconv"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/strings"
	st "github.com/alexfalkowski/go-service/v2/time"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/memorystore"
	"go.uber.org/fx"
)

// KeyFunc for the limiter.
type KeyFunc func(context.Context) meta.Value

var (
	// ErrMissingKey for limiter.
	ErrMissingKey = errors.New("limiter: missing key")

	keys = map[string]KeyFunc{}
)

// RegisterKey with name and fn. Last register wins.
func RegisterKey(name string, fn KeyFunc) {
	keys[name] = fn
}

// New limiter.
func New(lc fx.Lifecycle, cfg *Config) (*Limiter, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	k, ok := keys[cfg.Kind]
	if !ok {
		return nil, ErrMissingKey
	}

	// Memory store does not return an error.
	store, _ := memorystore.New(&memorystore.Config{
		Tokens:   cfg.Tokens,
		Interval: st.MustParseDuration(cfg.Interval),
	})

	limiter := &Limiter{
		store: store,
		key:   k,
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return limiter.Close(ctx)
		},
	})

	return limiter, nil
}

// Limiter holds a store with a key.
type Limiter struct {
	store limiter.Store
	key   KeyFunc
}

// Take from the store, returns if successful, info and error.
func (l *Limiter) Take(ctx context.Context) (bool, string, error) {
	tokens, remaining, _, ok, err := l.store.Take(ctx, l.key(ctx).Value())
	if err != nil {
		return false, "", err
	}

	header := strings.Concat(
		"limit=",
		strconv.FormatUint(tokens, 10),
		", remaining=",
		strconv.FormatUint(remaining, 10),
	)

	return ok, header, nil
}

// Close the limiter.
func (l *Limiter) Close(ctx context.Context) error {
	return l.store.Close(ctx)
}

package limiter

import (
	"context"
	"errors"
	"strconv"

	"github.com/alexfalkowski/go-service/meta"
	st "github.com/alexfalkowski/go-service/time"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/memorystore"
	"go.uber.org/fx"
)

type (
	// KeyFunc for the limiter.
	KeyFunc func(context.Context) meta.Valuer

	// Limiter limits the number of requests.
	Limiter interface {
		// Take a request.
		Take(ctx context.Context) (bool, string, error)

		// Close the limiter.
		Close(ctx context.Context) error
	}
)

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
//
//nolint:nilnil
func New(lc fx.Lifecycle, cfg *Config) (Limiter, error) {
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

	limiter := &Memory{
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

// Memory limiter.
type Memory struct {
	store limiter.Store
	key   KeyFunc
}

// Take from the store, returns if successful, info and error.
func (m *Memory) Take(ctx context.Context) (bool, string, error) {
	tokens, remaining, _, ok, err := m.store.Take(ctx, meta.ValueOrBlank(m.key(ctx)))
	if err != nil {
		return false, "", err
	}

	v := "limit=" + strconv.FormatUint(tokens, 10) + ", remaining=" + strconv.FormatUint(remaining, 10)

	return ok, v, nil
}

// Close the limiter.
func (m *Memory) Close(ctx context.Context) error {
	return m.store.Close(ctx)
}

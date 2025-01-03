package limiter

import (
	"context"
	"errors"
	"strconv"

	"github.com/alexfalkowski/go-service/meta"
	st "github.com/alexfalkowski/go-service/time"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/memorystore"
)

// KeyFunc for the limiter.
type KeyFunc func(context.Context) meta.Valuer

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
func New(cfg *Config) (limiter.Store, KeyFunc, error) {
	if !IsEnabled(cfg) {
		return nil, nil, nil
	}

	k, ok := keys[cfg.Kind]
	if !ok {
		return nil, nil, ErrMissingKey
	}

	store, err := memorystore.New(&memorystore.Config{Tokens: cfg.Tokens, Interval: st.MustParseDuration(cfg.Interval)})
	if err != nil {
		return nil, nil, err
	}

	return store, k, nil
}

// Take from the store, returns if successful, info and error.
func Take(ctx context.Context, store limiter.Store, key KeyFunc) (bool, string, error) {
	tokens, remaining, _, ok, err := store.Take(ctx, meta.ValueOrBlank(key(ctx)))
	if err != nil {
		return false, "", err
	}

	v := "limit=" + strconv.FormatUint(tokens, 10) + ", remaining=" + strconv.FormatUint(remaining, 10)

	return ok, v, nil
}

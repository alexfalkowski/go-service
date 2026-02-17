package limiter

import (
	"strconv"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/memorystore"
)

// KeyFunc returns the meta.Value used to key rate limits for ctx.
type KeyFunc func(context.Context) meta.Value

// KeyMap maps a kind string to the KeyFunc used to derive the limiter key.
type KeyMap map[string]KeyFunc

// NewKeyMap returns the default KeyMap used by the limiter.
//
// Supported default kinds are:
//   - "user-agent" -> meta.UserAgent
//   - "ip" -> meta.IPAddr
//   - "token" -> meta.Authorization
func NewKeyMap() KeyMap {
	return KeyMap{
		"user-agent": meta.UserAgent,
		"ip":         meta.IPAddr,
		"token":      meta.Authorization,
	}
}

// ErrMissingKey is returned when the configured key kind is not present in the KeyMap.
var ErrMissingKey = errors.New("limiter: missing key")

// NewLimiter constructs a Limiter using the configured key kind and interval/tokens settings.
//
// It returns ErrMissingKey when cfg.Kind is not present in keys.
func NewLimiter(lc di.Lifecycle, keys KeyMap, cfg *Config) (*Limiter, error) {
	k, ok := keys[cfg.Kind]
	if !ok {
		return nil, ErrMissingKey
	}

	interval := time.MustParseDuration(cfg.Interval)
	config := &memorystore.Config{
		Tokens:        cfg.Tokens,
		Interval:      interval,
		SweepMinTTL:   time.Hour,
		SweepInterval: time.Hour,
	}
	store, _ := memorystore.New(config)
	limiter := &Limiter{store: store, key: k}

	lc.Append(di.Hook{
		OnStop: func(ctx context.Context) error {
			return limiter.Close(ctx)
		},
	})

	return limiter, nil
}

// Limiter holds a store and a KeyFunc used to derive per-request limit keys.
type Limiter struct {
	store limiter.Store
	key   KeyFunc
}

// Take attempts to take a token for the key derived from ctx.
//
// It returns ok=false when the rate limit is exceeded. The returned header value is formatted as:
// "limit=<tokens>, remaining=<remaining>".
func (l *Limiter) Take(ctx context.Context) (bool, string, error) {
	tokens, remaining, _, ok, err := l.store.Take(ctx, l.key(ctx).Value())
	if err != nil {
		return false, strings.Empty, err
	}

	header := strings.Concat(
		"limit=",
		strconv.FormatUint(tokens, 10),
		", remaining=",
		strconv.FormatUint(remaining, 10),
	)
	return ok, header, nil
}

// Close closes the underlying store.
func (l *Limiter) Close(ctx context.Context) error {
	return l.store.Close(ctx)
}

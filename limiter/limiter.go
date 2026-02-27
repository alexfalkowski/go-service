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

// KeyFunc derives the metadata value used to key rate limits for ctx.
//
// The returned meta.Value is expected to yield a stable string via Value() that can be used as a
// per-request/per-actor limiter key (for example a user-agent, an IP address, or an authorization token).
type KeyFunc func(context.Context) meta.Value

// KeyMap maps a configured kind string to the KeyFunc used to derive the limiter key.
//
// It is typically constructed via NewKeyMap and passed to NewLimiter along with Config.Kind.
type KeyMap map[string]KeyFunc

// NewKeyMap returns the default KeyMap used by the limiter.
//
// Supported default kinds are:
//   - "user-agent": rate limit per User-Agent header (meta.UserAgent)
//   - "ip": rate limit per client IP address (meta.IPAddr)
//   - "token": rate limit per authorization token/header (meta.Authorization)
//
// These are simple defaults and may not be appropriate for all deployments. For example,
// User-Agent can be spoofed, and IP-based keys behave differently behind NATs/proxies unless
// meta.IPAddr is populated correctly by upstream middleware.
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
// NewLimiter selects a KeyFunc using cfg.Kind from the provided keys map. It then creates an
// in-memory limiter store configured with cfg.Tokens and cfg.Interval.
//
// Lifecycle behavior:
//   - OnStop: closes the underlying store via (*Limiter).Close.
//
// Errors:
//   - Returns ErrMissingKey when cfg.Kind is not present in keys.
//
// Notes:
//   - cfg.Interval is parsed using time.MustParseDuration and will panic if invalid.
//   - The underlying store constructor currently does not return an error (it is ignored).
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

// Limiter enforces rate limits using a store and a KeyFunc used to derive per-request keys.
//
// Limits are enforced per derived key string (see KeyFunc/KeyMap). This limiter uses an in-memory
// store, so limits are process-local and are not shared across replicas.
type Limiter struct {
	store limiter.Store
	key   KeyFunc
}

// Take attempts to take a token for the key derived from ctx.
//
// It delegates to the underlying limiter store using the derived key string:
//
//	l.key(ctx).Value()
//
// Return values:
//   - ok: false when the rate limit is exceeded for the derived key, true otherwise.
//   - header: a human-readable header value formatted as:
//     "limit=<tokens>, remaining=<remaining>".
//     The values represent the store-reported token limit and remaining tokens after this attempt.
//   - error: any error returned by the underlying store.
//
// Note: callers are responsible for deciding how to surface the header value (e.g. in HTTP response headers).
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

// Close closes the underlying store and releases any associated resources.
//
// This is typically invoked automatically via the lifecycle hook installed by NewLimiter.
func (l *Limiter) Close(ctx context.Context) error {
	return l.store.Close(ctx)
}

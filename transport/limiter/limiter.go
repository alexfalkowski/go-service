package limiter

import (
	"crypto/sha256"
	"encoding/hex"
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

// maxKeySize keeps normal limiter keys readable while preventing oversized
// request metadata from being retained verbatim as in-memory bucket names.
const maxKeySize = 256

// KeyFunc derives the metadata value used to key rate limits for ctx.
//
// The returned [meta.Value] is expected to yield a stable string via Value() that can be used as a
// per-request/per-actor limiter key (for example a user-agent, an IP address, a transport method, or a
// verified principal).
type KeyFunc func(context.Context) meta.Value

// KeyMap maps a configured kind string to the KeyFunc used to derive the limiter key.
//
// It is typically constructed via NewKeyMap and passed to NewLimiter along with [Config.Kind].
type KeyMap map[string]KeyFunc

// NewKeyMap returns the default KeyMap used by the limiter.
//
// Supported default kinds are:
//   - "user-agent": rate limit per User-Agent header ([meta.UserAgent])
//   - "ip": rate limit per client IP address ([meta.IPAddr])
//   - "user-id": rate limit per verified user/principal identifier ([meta.UserID])
//   - "service-method": rate limit per HTTP route/path or gRPC full method ([meta.ServiceMethod])
//
// These defaults are intended for controlled service-to-service traffic where user agents,
// forwarded IP headers, and authorization metadata are supplied by trusted clients or platform
// infrastructure. They are not sufficient as public-edge anti-abuse controls when clients can
// freely spoof headers; use trusted ingress, gateway, service-mesh, or post-auth identity limits
// for those boundaries.
func NewKeyMap() KeyMap {
	return KeyMap{
		"user-agent":     meta.UserAgent,
		"ip":             meta.IPAddr,
		"service-method": meta.ServiceMethod,
		"user-id":        meta.UserID,
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
//   - OnStop: closes the underlying store via [Limiter.Close].
//
// Errors:
//   - Returns ErrMissingKey when cfg.Kind is not present in keys or maps to a nil KeyFunc.
//
// Notes:
//   - cfg.Interval is used directly as a typed duration decoded from config.
//   - The underlying store constructor currently does not return an error (it is ignored).
func NewLimiter(lc di.Lifecycle, keys KeyMap, cfg *Config) (*Limiter, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	k, ok := keys[cfg.Kind]
	if !ok || k == nil {
		return nil, ErrMissingKey
	}

	config := &memorystore.Config{
		Tokens:   cfg.Tokens,
		Interval: cfg.Interval.Duration(),
		// Keep buckets at least as long as the configured limiter window so long
		// intervals are not purged and reset before the window completes.
		SweepMinTTL:   max(time.Hour.Duration(), cfg.Interval.Duration()),
		SweepInterval: time.Hour.Duration(),
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
	tokens, remaining, _, ok, err := l.store.Take(ctx, key(l.key(ctx)))
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

func key(value meta.Value) string {
	rawKey := value.Value()
	if strings.IsEmpty(rawKey) {
		return "<empty>"
	}
	if len(rawKey) <= maxKeySize {
		return rawKey
	}

	sum := sha256.Sum256([]byte(rawKey))
	return strings.Concat("sha256:", hex.EncodeToString(sum[:]))
}

package limiter

import (
	"math"
	"strconv"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/memorystore"
)

// ErrMissingKey is returned when the configured key kind is not present in the KeyMap.
var ErrMissingKey = errors.New("limiter: missing key")

// ErrIntervalTooLarge is returned when the configured interval cannot be used safely.
var ErrIntervalTooLarge = errors.New("limiter: interval too large")

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
//   - Returns ErrIntervalTooLarge when cfg.Interval would overflow internal key tracking TTLs.
//
// Notes:
//   - cfg.Interval is used directly as a typed duration decoded from config.
//     The underlying in-memory store applies its default 1s interval when cfg.Interval is 0.
//   - The underlying store constructor currently does not return an error (it is ignored).
func NewLimiter(lc di.Lifecycle, keyMap KeyMap, cfg *Config) (*Limiter, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	k, ok := keyMap[cfg.Kind]
	if !ok || k == nil {
		return nil, ErrMissingKey
	}

	sweepMinTTL := max(time.Hour.Duration(), cfg.Interval.Duration())
	sweepInterval := time.Hour.Duration()
	if sweepMinTTL > time.Duration(math.MaxInt64).Duration()-sweepInterval {
		return nil, ErrIntervalTooLarge
	}

	config := &memorystore.Config{
		Tokens:   cfg.Tokens,
		Interval: cfg.Interval.Duration(),
		// Keep buckets at least as long as the configured limiter window so long
		// intervals are not purged and reset before the window completes.
		SweepMinTTL:   sweepMinTTL,
		SweepInterval: sweepInterval,
	}
	store, _ := memorystore.New(config)
	limiter := &Limiter{
		store: store,
		key:   k,
		keys: &keys{
			values:  map[string]time.Time{},
			ttl:     time.Duration(sweepMinTTL + sweepInterval),
			maxKeys: cfg.GetMaxKeys(),
		},
	}

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
	keys  *keys
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
	decision, err := l.TakeDecision(ctx)
	return decision.Allowed(), decision.Header(), err
}

// TakeDecision attempts to take a token and returns the full limiter decision.
func (l *Limiter) TakeDecision(ctx context.Context) (Decision, error) {
	tokens, remaining, reset, ok, err := l.store.Take(ctx, l.keys.storeKey(l.key(ctx)))
	if err != nil {
		return Decision{}, err
	}

	header := strings.Concat(
		"limit=",
		strconv.FormatUint(tokens, 10),
		", remaining=",
		strconv.FormatUint(remaining, 10),
	)

	return Decision{
		allowed:    ok,
		header:     header,
		resetAfter: resetAfter(reset),
	}, nil
}

func resetAfter(reset uint64) time.Duration {
	now := uint64(time.Now().UnixNano())
	if reset <= now {
		return 0
	}

	//nolint:gosec // Reset delay is bounded by the validated limiter interval from the internal store.
	return time.Duration(reset - now)
}

// Close closes the underlying store and releases any associated resources.
//
// This is typically invoked automatically via the lifecycle hook installed by NewLimiter.
func (l *Limiter) Close(ctx context.Context) error {
	return l.store.Close(ctx)
}

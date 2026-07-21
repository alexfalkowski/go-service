package limiter

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-sync"
)

// maxKeySize keeps normal limiter keys readable while preventing oversized
// request metadata from being retained verbatim as in-memory bucket names.
const maxKeySize = 256

// overflowStoreKey is outside the storeKey namespace, which only emits empty:, raw:, and hash: prefixes.
const overflowStoreKey = "internal:max-keys-overflow"

// NewKeyMap returns the default KeyMap used by the limiter.
//
// Supported default kinds, in the usual order of preference, are:
//   - "user-id": rate limit per verified user/principal identifier ([meta.UserID])
//   - "transport-service-method": rate limit per transport-prefixed service method
//     ([meta.TransportServiceMethod])
//   - "service-method": rate limit per HTTP route/path or gRPC full method ([meta.ServiceMethod])
//   - "ip": rate limit per client IP address ([meta.IPAddr])
//   - "user-agent": rate limit per User-Agent header ([meta.UserAgent])
//
// These defaults are intended for controlled service-to-service traffic where user agents,
// forwarded IP headers, and authorization metadata are supplied by trusted clients or platform
// infrastructure. They are not sufficient as public-edge anti-abuse controls when clients can
// freely spoof headers; use trusted ingress, gateway, service-mesh, or post-auth identity limits
// for those boundaries.
func NewKeyMap() KeyMap {
	return KeyMap{
		"user-id":                  meta.UserID,
		"transport-service-method": meta.TransportServiceMethod,
		"service-method":           meta.ServiceMethod,
		"ip":                       meta.IPAddr,
		"user-agent":               meta.UserAgent,
	}
}

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

type keys struct {
	values        map[string]time.Time
	lastSweep     time.Time
	lock          sync.Mutex
	ttl           time.Duration
	sweepInterval time.Duration
	maxKeys       uint64
}

func (k *keys) storeKey(value meta.Value) string {
	key := storeKey(value)
	now := time.Now()

	k.lock.Lock()
	defer k.lock.Unlock()

	if _, ok := k.values[key]; ok {
		k.values[key] = now
		return key
	}

	if uint64(len(k.values)) < k.maxKeys {
		k.values[key] = now
		return key
	}

	// Match the store's sweep cadence so full-map expiry cleanup is amortized during key floods.
	if k.lastSweep.IsZero() || now.Sub(k.lastSweep) >= k.sweepInterval.Duration() {
		k.deleteExpired(now)
		k.lastSweep = now
	}

	if uint64(len(k.values)) < k.maxKeys {
		k.values[key] = now
		return key
	}

	return overflowStoreKey
}

func (k *keys) deleteExpired(now time.Time) {
	expiresBefore := now.Add(-k.ttl.Duration())
	for key, seenAt := range k.values {
		if seenAt.Before(expiresBefore) {
			delete(k.values, key)
		}
	}
}

func storeKey(value meta.Value) string {
	rawKey := value.Value()
	// Namespace every store key representation so caller-controlled raw values
	// cannot collide with internal empty sentinels or oversized-key hashes.
	if strings.IsEmpty(rawKey) {
		return "empty:"
	}
	if len(rawKey) <= maxKeySize {
		return strings.Concat("raw:", strconv.Itoa(len(rawKey)), ":", rawKey)
	}

	sum := sha256.Sum256([]byte(rawKey))
	return strings.Concat("hash:sha256:", hex.EncodeToString(sum[:]))
}

package limiter

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"

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

type keys struct {
	values  map[string]time.Time
	lock    sync.Mutex
	ttl     time.Duration
	maxKeys uint64
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

	k.deleteExpired(now)
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

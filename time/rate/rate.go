package rate

import (
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
	"golang.org/x/time/rate"
)

// New limiter for rate.
func New(cache *ristretto.Cache, every time.Duration, burst uint) *Limiter {
	return &Limiter{cache: cache, every: every, burst: burst}
}

// Limiter is a set of rate limiters for specific ids.
type Limiter struct {
	every time.Duration
	burst uint
	cache *ristretto.Cache
	lock  sync.RWMutex
}

// Get a limiter that is specified for that id.
func (u *Limiter) Get(id string) *rate.Limiter {
	u.lock.RLock()

	if r, ok := u.cache.Get(id); ok {
		u.lock.RUnlock()

		return r.(*rate.Limiter) // nolint:forcetypeassert
	}

	u.lock.RUnlock()

	u.lock.Lock()
	defer u.lock.Unlock()

	limiter := rate.NewLimiter(rate.Every(u.every), int(u.burst))

	u.cache.SetWithTTL(id, limiter, 0, time.Hour)

	return limiter
}

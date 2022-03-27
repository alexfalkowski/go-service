package rate

import (
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
	"golang.org/x/time/rate"
)

// Params for rate.
type Params struct {
	Every time.Duration
	Burst uint
	TTL   time.Duration
	Cache *ristretto.Cache
}

// New limiter for rate.
func New(params Params) *Limiter {
	return &Limiter{every: params.Every, burst: params.Burst, ttl: params.TTL, cache: params.Cache}
}

// Limiter is a set of rate limiters for specific ids.
type Limiter struct {
	every time.Duration
	burst uint
	ttl   time.Duration
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

	u.cache.SetWithTTL(id, limiter, 0, u.ttl)

	return limiter
}

package rate

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// New limiter for rate.
func New(every time.Duration, burst uint) *Limiter {
	return &Limiter{
		limiters: map[string]*rate.Limiter{},
		every:    every,
		burst:    burst,
	}
}

// Limiter is a set of rate limiters for specific ids.
type Limiter struct {
	every    time.Duration
	burst    uint
	limiters map[string]*rate.Limiter
	lock     sync.RWMutex
}

// Limiter that is specified for that id.
func (u *Limiter) Limiter(id string) *rate.Limiter {
	u.lock.RLock()

	if r, ok := u.limiters[id]; ok {
		u.lock.RUnlock()

		return r
	}

	u.lock.RUnlock()

	u.lock.Lock()
	defer u.lock.Unlock()

	u.limiters[id] = rate.NewLimiter(rate.Every(u.every), int(u.burst))

	return u.limiters[id]
}

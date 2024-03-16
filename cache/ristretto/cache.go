package ristretto

import (
	"time"
)

// Cache for ristretto.
type Cache interface {
	Get(key any) (any, bool)
	Set(key, value any, cost int64) bool
	SetWithTTL(key, value any, cost int64, ttl time.Duration) bool
	Del(key any)
	GetTTL(key any) (time.Duration, bool)
	Close()

	Hits() uint64
	Misses() uint64
}

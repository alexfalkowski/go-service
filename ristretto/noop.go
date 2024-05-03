package ristretto

import (
	"time"
)

// NewNoopCache for ristretto.
func NewNoopCache() *NoopCache {
	return &NoopCache{}
}

// NoopCache for ristretto.
type NoopCache struct{}

func (*NoopCache) Get(_ any) (any, bool) {
	return nil, false
}

func (*NoopCache) Set(_, _ any, _ int64) bool {
	return false
}

func (*NoopCache) SetWithTTL(_, _ any, _ int64, _ time.Duration) bool {
	return false
}

func (*NoopCache) Del(_ any) {
}

func (*NoopCache) GetTTL(_ any) (time.Duration, bool) {
	return 0, false
}

func (*NoopCache) Close() {
}

func (*NoopCache) Hits() uint64 {
	return 0
}

func (*NoopCache) Misses() uint64 {
	return 0
}

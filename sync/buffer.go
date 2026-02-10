package sync

import "github.com/alexfalkowski/go-sync"

// NewBufferPool is an alias for go-sync.NewBufferPool.
func NewBufferPool() *BufferPool {
	return sync.NewBufferPool()
}

// BufferPool is an alias for go-sync.BufferPool.
type BufferPool = sync.BufferPool

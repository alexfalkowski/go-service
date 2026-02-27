package sync

import "github.com/alexfalkowski/go-sync"

// NewBufferPool constructs a new BufferPool.
//
// This function forwards to github.com/alexfalkowski/go-sync.NewBufferPool and
// returns a *BufferPool (which is a type alias of the upstream implementation).
//
// Buffer pools are typically used to reduce allocations in hot paths that build
// byte payloads (for example encoders, compressors, and transports) by reusing
// temporary buffers across operations.
//
// Consult the upstream go-sync documentation for details of the pool’s API and
// usage expectations (for example how buffers are acquired/reset/released).
func NewBufferPool() *BufferPool {
	return sync.NewBufferPool()
}

// BufferPool provides pooled buffers to reduce allocations.
//
// BufferPool is a type alias of github.com/alexfalkowski/go-sync.BufferPool.
// Because it is an alias, its behavior, method set, and performance
// characteristics are those of the upstream implementation.
type BufferPool = sync.BufferPool

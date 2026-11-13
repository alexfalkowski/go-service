package sync

import (
	"github.com/alexfalkowski/go-signal/v2"
	"github.com/alexfalkowski/go-sync"
)

// Module wires the upstream go-sync buffer pool into [github.com/alexfalkowski/go-signal/v2].
//
// Including this module in a go-signal application provides a single shared buffer
// pool dependency, [github.com/alexfalkowski/go-sync.BufferPool], constructed
// via [github.com/alexfalkowski/go-sync.NewBufferPool].
//
// This is commonly used to reduce allocations across components that build or
// transform byte payloads, such as encoders, compressors, caches, and HTTP
// helpers.
//
// Module only provides the buffer pool. It does not register lifecycle hooks,
// configure the pool, or re-export other upstream synchronization helpers.
//
// Standard go-service module bundles already include this module through their
// shared library wiring, so custom compositions should include it only when they
// are not already using one of those bundles.
//
// The provided type and its lifecycle/usage contract are defined by the
// upstream [github.com/alexfalkowski/go-sync] package.
var Module = signal.Module(
	signal.Constructor(sync.NewBufferPool),
)

package sync

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-sync"
)

// Module wires the upstream go-sync buffer pool into Fx.
//
// Including this module in an Fx application provides a single shared buffer
// pool instance constructed via github.com/alexfalkowski/go-sync.NewBufferPool.
//
// This is commonly used to reduce allocations across components that build or
// transform byte payloads, such as encoders, compressors, caches, and HTTP
// helpers.
//
// The provided type and its lifecycle/usage contract are defined by the
// upstream github.com/alexfalkowski/go-sync package.
var Module = di.Module(
	di.Constructor(sync.NewBufferPool),
)

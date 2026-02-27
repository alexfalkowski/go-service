package sync

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-sync"
)

// Module wires a shared buffer pool into Fx.
//
// Including this module in an Fx application provides a single, shared
// *BufferPool instance constructed via github.com/alexfalkowski/go-sync.NewBufferPool.
//
// This is commonly used to reduce allocations across components that build or
// transform byte payloads (for example encoders, compressors, and transports)
// by reusing temporary buffers.
//
// Note: the concrete BufferPool type and its lifecycle/usage contract are
// defined by the upstream go-sync package.
var Module = di.Module(
	di.Constructor(sync.NewBufferPool),
)

package sync

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-sync"
)

// Module wires the buffer pool into Fx.
var Module = di.Module(
	di.Constructor(sync.NewBufferPool),
)

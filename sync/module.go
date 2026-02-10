package sync

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-sync"
)

// Module for fx.
var Module = di.Module(
	di.Constructor(sync.NewBufferPool),
)

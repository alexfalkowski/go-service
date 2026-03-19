package test

import "github.com/alexfalkowski/go-service/v2/sync"

// Pool is the shared buffer pool used by transport and encoding test helpers.
var Pool = sync.NewBufferPool()

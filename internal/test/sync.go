package test

import "github.com/alexfalkowski/go-sync"

// Pool is the shared buffer pool used by transport and encoding test helpers.
var Pool = sync.NewBufferPool()

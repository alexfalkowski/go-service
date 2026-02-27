package sync

import "github.com/alexfalkowski/go-sync"

// Mutex is a mutual exclusion lock.
//
// Mutex is a type alias of github.com/alexfalkowski/go-sync.Mutex (which in turn
// is expected to mirror the semantics of the standard library sync.Mutex).
//
// Use Mutex to protect shared state that must not be accessed concurrently. The
// zero value is usable without initialization.
type Mutex = sync.Mutex

// RWMutex is a reader/writer mutual exclusion lock.
//
// RWMutex is a type alias of github.com/alexfalkowski/go-sync.RWMutex (which in
// turn is expected to mirror the semantics of the standard library sync.RWMutex).
//
// A RWMutex allows concurrent readers or a single writer. Use RWMutex when reads
// are frequent and can safely proceed in parallel, while writes must be exclusive.
// The zero value is usable without initialization.
type RWMutex = sync.RWMutex

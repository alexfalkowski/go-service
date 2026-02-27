package sync

import "github.com/alexfalkowski/go-sync"

// NewMap constructs a new concurrent Map.
// //
// // This function forwards to github.com/alexfalkowski/go-sync.NewMap and returns
// // a Map[K, V] (which is a type alias of the upstream implementation).
// //
// // Use Map when you need a concurrency-safe key/value store without manually
// // managing locks. The exact method set and concurrency guarantees are defined
// // by the upstream go-sync package.
func NewMap[K comparable, V any]() Map[K, V] {
	return sync.NewMap[K, V]()
}

// Map is a generic concurrency-safe map keyed by K with values of type V.
//
// Map is a type alias of github.com/alexfalkowski/go-sync.Map. Because it is an
// alias, its behavior, method set, and performance characteristics are those of
// the upstream implementation.
//
// If you need semantics beyond what the upstream Map provides (for example
// custom eviction policies or strong iteration guarantees), use a different data
// structure or explicit locking with Mutex/RWMutex.
type Map[K comparable, V any] = sync.Map[K, V]

package sync

import "github.com/alexfalkowski/go-sync"

// NewMap is an alias for go-sync.NewMap.
func NewMap[K comparable, V any]() Map[K, V] {
	return sync.NewMap[K, V]()
}

// Map is an alias for go-sync.Map.
type Map[K comparable, V any] = sync.Map[K, V]

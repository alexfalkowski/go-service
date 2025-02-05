package sync

import (
	"sync"

	"github.com/alexfalkowski/go-service/types/slices"
)

// NewArray fir sync.
func NewArray[T any]() *Array[T] {
	return &Array[T]{arr: make([]*T, 0)}
}

// Array allows an array of T to safely add and read items.
type Array[T any] struct {
	arr []*T
	mux sync.RWMutex
}

// Add an elemment.
func (a *Array[T]) Add(elem *T) {
	a.mux.Lock()
	defer a.mux.Unlock()

	a.arr = slices.Append(a.arr, elem)
}

// Elements of the array.
func (a *Array[T]) Elements() []*T {
	a.mux.RLock()
	defer a.mux.RUnlock()

	return a.arr
}

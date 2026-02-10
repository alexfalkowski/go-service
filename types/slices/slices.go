package slices

import (
	"slices"

	"github.com/alexfalkowski/go-service/v2/types/structs"
)

// AppendNotZero appends elems to slice, skipping elements that are the zero value for T.
func AppendNotZero[T comparable](slice []T, elems ...T) []T {
	for _, elem := range elems {
		if structs.IsZero(elem) {
			continue
		}
		slice = append(slice, elem)
	}
	return slice
}

// AppendNotNil appends elems to slice, skipping nil elements.
func AppendNotNil[T any](slice []*T, elems ...*T) []*T {
	for _, elem := range elems {
		if structs.IsNil(elem) {
			continue
		}
		slice = append(slice, elem)
	}
	return slice
}

// ElemFunc returns the first element in slice that matches f and whether it was found.
//
// It is equivalent to using `slices.IndexFunc` and then indexing, but returns the element directly.
func ElemFunc[T any](slice []*T, f func(*T) bool) (*T, bool) {
	index := slices.IndexFunc(slice, f)
	if index == -1 {
		return nil, false
	}

	return slice[index], true
}

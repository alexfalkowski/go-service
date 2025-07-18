package slices

import (
	"slices"

	"github.com/alexfalkowski/go-service/v2/types/structs"
)

// AppendNotZero elements to the slice, only if the element is not zero.
func AppendNotZero[T comparable](slice []T, elems ...T) []T {
	for _, elem := range elems {
		if structs.IsZero(elem) {
			continue
		}
		slice = append(slice, elem)
	}
	return slice
}

// AppendNotNil elements to the slice, only if the element is not nil.
func AppendNotNil[T any](slice []*T, elems ...*T) []*T {
	for _, elem := range elems {
		if structs.IsNil(elem) {
			continue
		}
		slice = append(slice, elem)
	}
	return slice
}

// ElemFunc is just slices.IndexFunc, returning the element and if it was found.
func ElemFunc[T any](slice []*T, f func(*T) bool) (*T, bool) {
	index := slices.IndexFunc(slice, f)
	if index == -1 {
		return nil, false
	}

	return slice[index], true
}

package slices

import (
	"github.com/alexfalkowski/go-service/structs"
)

// AppendNotZero elements to the slice, only if the element is not zero.
func AppendNotZero[T comparable](slice []*T, elems ...*T) []*T {
	return Append(slice, structs.IsZero, elems...)
}

// AppendNotNil elements to the slice, only if the element is not nil.
func AppendNotNil[T any](slice []*T, elems ...*T) []*T {
	return Append(slice, func(t *T) bool { return t == nil }, elems...)
}

// Append elements to the slice, only if the element is not zero.
func Append[T any](slice []*T, eq func(*T) bool, elems ...*T) []*T {
	for _, elem := range elems {
		if eq(elem) {
			continue
		}

		slice = append(slice, elem)
	}

	return slice
}

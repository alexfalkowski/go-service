package slices

import (
	"github.com/alexfalkowski/go-service/structs"
)

// Append elements to the slice, only if the element is not zero.
func Append[T comparable](slice []*T, elems ...*T) []*T {
	for _, elem := range elems {
		if structs.IsZero[T](elem) {
			continue
		}

		slice = append(slice, elem)
	}

	return slice
}

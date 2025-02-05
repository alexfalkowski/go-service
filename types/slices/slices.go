package slices

import (
	"github.com/alexfalkowski/go-service/types/structs"
)

// Append elements to the slice, only if the element is not nil.
func Append[T any](slice []*T, elems ...*T) []*T {
	for _, elem := range elems {
		if structs.IsNil(elem) {
			continue
		}

		slice = append(slice, elem)
	}

	return slice
}

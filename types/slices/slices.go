package slices

import (
	"slices"

	"github.com/alexfalkowski/go-service/v2/types/structs"
)

// AppendNotZero appends elems to slice, skipping elements that are the zero value for T.
//
// “Zero” is determined using structs.IsZero, which compares the element against the
// type’s zero value using == (therefore T must be comparable).
//
// This helper preserves the relative order of appended elements and returns the
// resulting slice.
//
// Example:
//
//	var out []string
//	out = slices.AppendNotZero(out, "", "a", "", "b")
//	// out == []string{"a", "b"}
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
//
// This helper is useful when building slices of optional pointer values where nil
// indicates “not provided”. It preserves the relative order of appended elements
// and returns the resulting slice.
//
// Example:
//
//	var out []*int
//	var a = 1
//	out = slices.AppendNotNil(out, nil, &a, nil)
//	// out contains only &a
func AppendNotNil[T any](slice []*T, elems ...*T) []*T {
	for _, elem := range elems {
		if structs.IsNil(elem) {
			continue
		}
		slice = append(slice, elem)
	}
	return slice
}

// ElemFunc returns the first element in slice that matches f, along with a boolean
// indicating whether a match was found.
//
// This helper is equivalent to:
//   - computing slices.IndexFunc(slice, f), and
//   - returning slice[index] when index != -1,
//
// but returns the element directly instead of an index.
//
// If no element matches, ElemFunc returns (nil, false).
func ElemFunc[T any](slice []*T, f func(*T) bool) (*T, bool) {
	index := slices.IndexFunc(slice, f)
	if index == -1 {
		return nil, false
	}

	return slice[index], true
}

package slices

import (
	"cmp"
	"iter"
	"slices"

	"github.com/alexfalkowski/go-service/v2/reflect"
)

// AppendNotZero appends elems to slice, skipping nil and zero elements.
//
// Zero values are determined using [reflect.IsZero], which supports
// non-comparable values such as slices, maps, funcs, and structs containing them.
// Nil slices and maps are skipped; non-nil empty slices and maps are appended.
//
// This helper preserves the relative order of appended elements and returns the
// resulting slice.
//
// Example:
//
//	var out []string
//	out = slices.AppendNotZero(out, "", "a", "", "b")
//	// out == []string{"a", "b"}
func AppendNotZero[T any](slice []T, elems ...T) []T {
	for _, elem := range elems {
		if reflect.IsZero(elem) {
			continue
		}
		slice = append(slice, elem)
	}
	return slice
}

// Backward returns an iterator over index-value pairs in slice, traversing it backward.
//
// It is a thin wrapper around the standard library [slices.Backward].
func Backward[S ~[]E, E any](slice S) iter.Seq2[int, E] {
	return slices.Backward(slice)
}

// ElemFunc returns the first element in slice that matches f, along with a boolean
// indicating whether a match was found.
//
// This helper is equivalent to:
//   - computing [slices.IndexFunc](slice, f), and
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

// Clip removes unused capacity from slice.
//
// It is a thin wrapper around the standard library [slices.Clip]. Use it when a
// returned slice may share a backing array with neighboring values and future
// appends should allocate instead of modifying that shared backing array.
func Clip[T any](slice []T) []T {
	return slices.Clip(slice)
}

// Clone returns a copy of slice.
//
// It is a thin wrapper around the standard library [slices.Clone].
func Clone[S ~[]E, E any](slice S) S {
	return slices.Clone(slice)
}

// Collect collects values from seq into a new slice.
//
// It is a thin wrapper around the standard library [slices.Collect].
func Collect[E any](seq iter.Seq[E]) []E {
	return slices.Collect(seq)
}

// Sort sorts a slice of any ordered type in ascending order.
//
// It is a thin wrapper around the standard library [slices.Sort].
func Sort[S ~[]E, E cmp.Ordered](slice S) {
	slices.Sort(slice)
}

// Contains reports whether v is present in slice.
//
// It is a thin wrapper around the standard library [slices.Contains].
func Contains[S ~[]E, E comparable](slice S, v E) bool {
	return slices.Contains(slice, v)
}

// ContainsFunc reports whether at least one element of slice satisfies f.
//
// It is a thin wrapper around the standard library [slices.ContainsFunc].
func ContainsFunc[S ~[]E, E any](slice S, f func(E) bool) bool {
	return slices.ContainsFunc(slice, f)
}

// DeleteFunc removes any elements from slice for which del returns true.
//
// It is a thin wrapper around the standard library [slices.DeleteFunc].
func DeleteFunc[S ~[]E, E any](slice S, del func(E) bool) S {
	return slices.DeleteFunc(slice, del)
}

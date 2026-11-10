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

// All wraps [slices.All] to keep iteration on the project-owned surface.
func All[S ~[]E, E any](slice S) iter.Seq2[int, E] {
	return slices.All(slice)
}

// AppendSeq wraps [slices.AppendSeq] to keep appending on the project-owned surface.
func AppendSeq[S ~[]E, E any](slice S, seq iter.Seq[E]) S {
	return slices.AppendSeq(slice, seq)
}

// BinarySearch wraps [slices.BinarySearch] to keep searching on the project-owned surface.
func BinarySearch[S ~[]E, E cmp.Ordered](slice S, target E) (int, bool) {
	return slices.BinarySearch(slice, target)
}

// BinarySearchFunc wraps [slices.BinarySearchFunc] to keep searching on the project-owned surface.
func BinarySearchFunc[S ~[]E, E, T any](slice S, target T, compare func(E, T) int) (int, bool) {
	return slices.BinarySearchFunc(slice, target, compare)
}

// Chunk wraps [slices.Chunk] to keep iteration on the project-owned surface.
func Chunk[S ~[]E, E any](slice S, n int) iter.Seq[S] {
	return slices.Chunk(slice, n)
}

// Compact wraps [slices.Compact] to keep compaction on the project-owned surface.
func Compact[S ~[]E, E comparable](slice S) S {
	return slices.Compact(slice)
}

// CompactFunc wraps [slices.CompactFunc] to keep compaction on the project-owned surface.
func CompactFunc[S ~[]E, E any](slice S, equal func(E, E) bool) S {
	return slices.CompactFunc(slice, equal)
}

// Compare wraps [slices.Compare] to keep comparison on the project-owned surface.
func Compare[S ~[]E, E cmp.Ordered](first, second S) int {
	return slices.Compare(first, second)
}

// CompareFunc wraps [slices.CompareFunc] to keep comparison on the project-owned surface.
func CompareFunc[S1 ~[]E1, S2 ~[]E2, E1, E2 any](first S1, second S2, compare func(E1, E2) int) int {
	return slices.CompareFunc(first, second, compare)
}

// Concat wraps [slices.Concat] to keep concatenation on the project-owned surface.
func Concat[S ~[]E, E any](all ...S) S {
	return slices.Concat(all...)
}

// Delete wraps [slices.Delete] to keep deletion on the project-owned surface.
func Delete[S ~[]E, E any](slice S, i, j int) S {
	return slices.Delete(slice, i, j)
}

// Equal wraps [slices.Equal] to keep equality checks on the project-owned surface.
func Equal[S ~[]E, E comparable](first, second S) bool {
	return slices.Equal(first, second)
}

// EqualFunc wraps [slices.EqualFunc] to keep equality checks on the project-owned surface.
func EqualFunc[S1 ~[]E1, S2 ~[]E2, E1, E2 any](first S1, second S2, equal func(E1, E2) bool) bool {
	return slices.EqualFunc(first, second, equal)
}

// Grow wraps [slices.Grow] to keep capacity management on the project-owned surface.
func Grow[S ~[]E, E any](slice S, n int) S {
	return slices.Grow(slice, n)
}

// Index wraps [slices.Index] to keep searching on the project-owned surface.
func Index[S ~[]E, E comparable](slice S, value E) int {
	return slices.Index(slice, value)
}

// IndexFunc wraps [slices.IndexFunc] to keep searching on the project-owned surface.
func IndexFunc[S ~[]E, E any](slice S, f func(E) bool) int {
	return slices.IndexFunc(slice, f)
}

// Insert wraps [slices.Insert] to keep insertion on the project-owned surface.
func Insert[S ~[]E, E any](slice S, i int, values ...E) S {
	return slices.Insert(slice, i, values...)
}

// IsSorted wraps [slices.IsSorted] to keep sort checks on the project-owned surface.
func IsSorted[S ~[]E, E cmp.Ordered](slice S) bool {
	return slices.IsSorted(slice)
}

// IsSortedFunc wraps [slices.IsSortedFunc] to keep sort checks on the project-owned surface.
func IsSortedFunc[S ~[]E, E any](slice S, compare func(E, E) int) bool {
	return slices.IsSortedFunc(slice, compare)
}

// Max wraps [slices.Max] to keep extrema selection on the project-owned surface.
func Max[S ~[]E, E cmp.Ordered](slice S) E {
	return slices.Max(slice)
}

// MaxFunc wraps [slices.MaxFunc] to keep extrema selection on the project-owned surface.
func MaxFunc[S ~[]E, E any](slice S, compare func(E, E) int) E {
	return slices.MaxFunc(slice, compare)
}

// Min wraps [slices.Min] to keep extrema selection on the project-owned surface.
func Min[S ~[]E, E cmp.Ordered](slice S) E {
	return slices.Min(slice)
}

// MinFunc wraps [slices.MinFunc] to keep extrema selection on the project-owned surface.
func MinFunc[S ~[]E, E any](slice S, compare func(E, E) int) E {
	return slices.MinFunc(slice, compare)
}

// Repeat wraps [slices.Repeat] to keep repetition on the project-owned surface.
func Repeat[S ~[]E, E any](slice S, count int) S {
	return slices.Repeat(slice, count)
}

// Replace wraps [slices.Replace] to keep replacement on the project-owned surface.
func Replace[S ~[]E, E any](slice S, i, j int, values ...E) S {
	return slices.Replace(slice, i, j, values...)
}

// Reverse wraps [slices.Reverse] to keep reversal on the project-owned surface.
func Reverse[S ~[]E, E any](slice S) {
	slices.Reverse(slice)
}

// SortFunc wraps [slices.SortFunc] to keep sorting on the project-owned surface.
func SortFunc[S ~[]E, E any](slice S, compare func(E, E) int) {
	slices.SortFunc(slice, compare)
}

// SortStableFunc wraps [slices.SortStableFunc] to keep stable sorting on the project-owned surface.
func SortStableFunc[S ~[]E, E any](slice S, compare func(E, E) int) {
	slices.SortStableFunc(slice, compare)
}

// Sorted wraps [slices.Sorted] to keep iterator sorting on the project-owned surface.
func Sorted[E cmp.Ordered](seq iter.Seq[E]) []E {
	return slices.Sorted(seq)
}

// SortedFunc wraps [slices.SortedFunc] to keep iterator sorting on the project-owned surface.
func SortedFunc[E any](seq iter.Seq[E], compare func(E, E) int) []E {
	return slices.SortedFunc(seq, compare)
}

// SortedStableFunc wraps [slices.SortedStableFunc] to keep stable iterator sorting on the project-owned surface.
func SortedStableFunc[E any](seq iter.Seq[E], compare func(E, E) int) []E {
	return slices.SortedStableFunc(seq, compare)
}

// Values wraps [slices.Values] to keep iteration on the project-owned surface.
func Values[S ~[]E, E any](slice S) iter.Seq[E] {
	return slices.Values(slice)
}

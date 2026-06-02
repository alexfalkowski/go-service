// Package slices provides small, generic helpers for working with slices.
//
// The helpers in this package are intentionally lightweight and are designed to
// reduce repetitive boilerplate in common slice-manipulation patterns used across
// go-service.
//
// This package complements the Go standard library [slices] package rather than
// replacing it. Prefer the standard library first when it provides what you need,
// and use these helpers when they improve readability or avoid repeated nil/zero
// checks.
//
// # Conditional append helpers
//
// [AppendNotZero] appends elements to a slice while skipping elements that are
// nil or the zero value of the element type.
//
// This is useful when building option lists, attribute sets, or filter lists where
// zero values should be treated as "not provided".
// The helper preserves the relative order of the elements it appends and returns
// the resulting slice.
//
// # Searching helpers
//
// [ElemFunc] returns the first pointer element in a slice that matches a
// predicate, along with a boolean indicating whether a match was found.
//
// It is a convenience wrapper around [slices.IndexFunc] that returns the element
// directly instead of an index.
//
// # Capacity helpers
//
// [Clip] removes unused capacity from a slice. Use it when returning a slice
// that may share a backing array with neighboring values and future appends
// should not mutate that shared backing array.
//
// # Notes
//
// These helpers are designed to be predictable and side-effect free aside from
// modifying the returned slice as normal for append-style operations. They do not
// mutate elements beyond appending them. "Zero" is defined by [reflect.IsZero], so
// nil slices and maps are skipped, while non-nil empty slices and maps are
// appended.
package slices

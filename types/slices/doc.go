// Package slices provides small, generic helpers for working with slices.
//
// The helpers in this package are intentionally lightweight and are designed to
// reduce repetitive boilerplate in common slice-manipulation patterns used across
// go-service.
//
// This package complements the Go standard library `slices` package rather than
// replacing it. Prefer the standard library first when it provides what you need,
// and use these helpers when they improve readability or avoid repeated nil/zero
// checks.
//
// # Conditional append helpers
//
// AppendNotZero appends elements to a slice while skipping elements that are the
// zero value of the element type.
//
// This is useful when building option lists, attribute sets, or filter lists where
// zero values should be treated as “not provided”.
//
// AppendNotNil appends pointer elements to a slice while skipping nil pointers.
//
// This is useful when building slices of optional pointers where nil indicates the
// absence of a value.
//
// Both helpers preserve the relative order of the elements they append and return
// the resulting slice.
//
// # Searching helpers
//
// ElemFunc returns the first pointer element in a slice that matches a predicate,
// along with a boolean indicating whether a match was found.
//
// It is a convenience wrapper around `slices.IndexFunc` that returns the element
// directly instead of an index.
//
// # Notes
//
// These helpers are designed to be predictable and side-effect free aside from
// modifying the returned slice as normal for append-style operations. They do not
// mutate elements (beyond appending) and they do not perform deep equality checks;
// “zero” is defined by Go’s == comparison against the zero value for comparable
// types.
package slices

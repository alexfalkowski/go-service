// Package structs provides small, generic helpers for working with pointers and
// zero values.
//
// Despite the package name, these helpers are not limited to struct types. They
// are intended to reduce boilerplate around common “is this value present?”
// checks that appear throughout go-service configuration and wiring.
//
// # Concepts: nil, zero, and empty
//
// This package distinguishes three closely related concepts:
//
//   - nil: the pointer itself is nil (no value is present).
//
//   - zero: a non-pointer value equals its type’s zero value (for example 0 for
//     integers, "" for strings, false for bools, an empty struct value, etc.).
//     In this package, “zero” is determined via == comparison against the zero
//     value of the type.
//
//   - empty: a pointer is considered “empty” when either it is nil or it points
//     to the zero value of its element type.
//
// # Functions
//
// IsNil reports whether a pointer is nil.
//
// IsZero reports whether a value equals the zero value for its type. Because it
// relies on ==, the type parameter must be comparable.
//
// IsEmpty reports whether a pointer is nil or points to a zero value.
//
// # Examples
//
//	nilPtr := (*int)(nil)
//	_ = structs.IsNil(nilPtr)   // true
//	_ = structs.IsEmpty(nilPtr) // true
//
//	v := 0
//	p := &v
//	_ = structs.IsNil(p)    // false
//	_ = structs.IsZero(v)   // true
//	_ = structs.IsEmpty(p)  // true (points to zero)
//
//	w := 42
//	q := &w
//	_ = structs.IsEmpty(q) // false
//
// # Notes and limitations
//
//   - IsZero and IsEmpty require comparable element types (T comparable). If you
//     need “zero” semantics for non-comparable types (for example slices, maps,
//     or structs containing non-comparable fields), you must define your own
//     notion of emptiness or use reflection-based checks (with the usual
//     tradeoffs).
//
//   - For pointer types, IsEmpty treats “pointer to zero” as empty. This is
//     useful for optional configuration values where a pointer may be present
//     but carry the default/zero value and you want to treat that as “unset”.
//     If you need to distinguish “unset (nil)” from “explicitly set to zero”,
//     do not use IsEmpty; check nil explicitly and handle the pointed-to value
//     separately.
package structs

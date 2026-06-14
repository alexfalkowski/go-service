// Package strings provides small string helpers and a curated set of aliases for
// the Go standard library strings package.
//
// The intent of this package is to:
//
//   - Keep go-service code depending on go-service packages consistently, while
//     still delegating to the standard library for core string operations.
//   - Provide a few convenience helpers that are frequently needed across the
//     repository (for example emptiness checks and small splitting helpers).
//   - Expose carefully documented unsafe utilities (see [Bytes]) for performance-
//     sensitive paths where avoiding allocations is important.
//
// Most functions in this package are thin wrappers around the corresponding
// functions in the standard library strings package (for example [Contains],
// [HasPrefix], [ReplaceAll], [TrimSpace]) or the strcase package (for example
// [ToSnake], [ToLowerCamel], [ToDelimited]). These wrappers do not change
// semantics.
//
// # Convenience helpers
//
// This package also provides helpers not present in the standard library:
//
//   - [IsEmpty] / [IsAnyEmpty]: common emptiness checks. [IsAnyEmpty]
//     returns false when called without values.
//
//   - [Join] and [Concat]: Join takes the separator first and then variadic
//     strings, so callers can avoid allocating a slice at the call site. With no
//     strings it returns [Empty]. Concat concatenates without a separator.
//
//   - [CutColon]: a small helper that splits on the first ":" using [strings.Cut],
//     returning the part before and after. If ":" is not present, "before" is
//     the original input and "after" is empty. Callers that need the found flag
//     should use [Cut] with ":" directly.
//
//   - [ToSnake] / [ToLowerCamel] / [ToDelimited]: common string case conversion
//     helpers delegated to strcase.
//
// # Constants
//
// [Empty] and [Space] are provided as named constants for readability and reuse.
//
// # Unsafe conversions
//
// [Bytes] converts a string to a byte slice without allocation by using unsafe.
//
// Important: the returned []byte aliases the same memory as the input string.
//
//   - Treat the returned slice as read-only. Writing to it results in undefined
//     behavior.
//
//   - Do not retain the returned slice beyond the lifetime of the original
//     string value. In particular, do not store it in long-lived structures or
//     return it when the string was derived from a transient buffer.
//
// For safe conversions, use the built-in conversion []byte(s), which allocates
// a new byte slice.
package strings

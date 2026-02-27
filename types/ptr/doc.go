// Package ptr provides small, generic helpers for working with pointers.
//
// This package is intentionally tiny. It exists to reduce repetitive boilerplate
// when you need a pointer to a value (for example for optional configuration
// fields, APIs that take *T, or tests).
//
// # Helpers
//
// Zero returns a pointer to the zero value of T:
//
//	var p *int = ptr.Zero[int]() // points to 0
//
// Value returns a pointer to the provided value:
//
//	p := ptr.Value("hello") // *string
//
// # Semantics and pitfalls
//
// These helpers allocate storage for the pointed-to value and return a pointer to
// that storage. Each call returns a distinct pointer.
//
// For example, two calls produce two different pointers even if the values are equal:
//
//	a := ptr.Zero[int]()
//	b := ptr.Zero[int]()
//	// a != b, but *a == *b
//
// # When not to use this package
//
// If you already have an addressable variable, taking its address is typically
// clearer and avoids an extra helper call:
//
//	v := 42
//	p := &v
//
// Similarly, for composite literals you can take the address directly:
//
//	p := &MyStruct{Field: "x"}
//
// This package is best used when you need a pointer but don't have an addressable
// value at the call site or you want a concise, generic helper in tests and
// configuration wiring.
package ptr

package ptr

// Zero returns a pointer to the zero value of T.
//
// This helper allocates storage for a zero value of T and returns a pointer to it.
// Each call returns a distinct pointer.
//
// Example:
//
//	p := ptr.Zero[int]() // points to 0
func Zero[T any]() *T {
	var t T
	return &t
}

// Value returns a pointer to t.
//
// This helper allocates storage for the provided value and returns a pointer to it.
// Each call returns a distinct pointer, even if the values are equal.
//
// Example:
//
//	p := ptr.Value("hello") // *string pointing to "hello"
func Value[T any](t T) *T {
	return new(t)
}

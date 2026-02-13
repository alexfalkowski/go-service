package ptr

// Zero returns a pointer to the zero value of T.
func Zero[T any]() *T {
	var t T
	return &t
}

// Value returns a pointer to t.
func Value[T any](t T) *T {
	return new(t)
}

package ptr

// Empty pointer of *T.
func Empty[T any]() *T {
	var t T

	return &t
}

// Value pointer of value of t.
func Value[T any](t T) *T {
	return &t
}

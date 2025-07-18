package ptr

// Zero pointer of T.
func Zero[T any]() *T {
	var t T
	return &t
}

// Value pointer from value of t.
func Value[T any](t T) *T {
	return &t
}

package ptr

// Zero returns a pointer to the zero value of T.
//
// This helper allocates storage for a zero value of T and returns a pointer to it.
// Each call returns a distinct pointer.
func Zero[T any]() *T {
	var t T
	return &t
}

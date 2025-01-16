package types

// Pointer returns *T.
func Pointer[T any]() *T {
	var t T

	return &t
}

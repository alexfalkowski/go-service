package structs

// Zero value for the type.
func Zero[T any]() T {
	var zero T

	return zero
}

// IsZero for a specific type.
func IsZero[T comparable](value *T) bool {
	if value == nil {
		return true
	}

	z := Zero[T]()

	return z == *value
}

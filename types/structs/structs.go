package structs

// IsNil for a specific type.
func IsNil[T any](value *T) bool {
	return value == nil
}

// Zero value for the type.
func Zero[T any]() T {
	var zero T

	return zero
}

// IsZero for a specific type.
func IsZero[T comparable](value *T) bool {
	if IsNil(value) {
		return true
	}

	z := Zero[T]()

	return z == *value
}

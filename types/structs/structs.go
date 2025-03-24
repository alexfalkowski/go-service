package structs

// IsEmpty checks if T is nil or zero.
func IsEmpty[T comparable](value *T) bool {
	return IsNil(value) || IsZero(*value)
}

// IsNil for a specific type.
func IsNil[T any](value *T) bool {
	return value == nil
}

// IsZero for a specific type.
func IsZero[T comparable](value T) bool {
	var zero T

	return value == zero
}

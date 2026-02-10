package structs

// IsEmpty reports whether value is nil or points to the zero value of T.
func IsEmpty[T comparable](value *T) bool {
	return IsNil(value) || IsZero(*value)
}

// IsNil reports whether value is nil.
func IsNil[T any](value *T) bool {
	return value == nil
}

// IsZero reports whether value is the zero value of T.
func IsZero[T comparable](value T) bool {
	var zero T
	return value == zero
}

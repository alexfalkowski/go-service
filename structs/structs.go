package structs

// IsEmpty reports whether value is nil or points to the zero value of T.
//
// Nil pointers and pointers to zero values both return true. If "unset" and
// "explicitly set to zero" have different meanings, use [IsNil] first and then
// inspect the pointed-to value separately.
//
// Because it compares the pointed-to value with the zero value, T must be
// comparable.
func IsEmpty[T comparable](value *T) bool {
	return IsNil(value) || IsZero(*value)
}

// IsNil reports whether value is nil.
//
// This is a small helper that improves readability at call sites where a pointer
// represents an optional value.
func IsNil[T any](value *T) bool {
	return value == nil
}

// IsZero reports whether value equals the zero value of T.
//
// Because it uses == comparison against the zero value, T must be comparable.
// If you need "zero" semantics for non-comparable types (for example slices, maps,
// or structs containing non-comparable fields), define your own emptiness check.
func IsZero[T comparable](value T) bool {
	var zero T
	return value == zero
}

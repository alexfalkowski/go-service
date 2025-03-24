package structs

import (
	"context"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Register the validator.
func Register(v *validator.Validate) {
	validate = v
}

// Validate the struct.
func Validate[T any](t *T) error {
	return validate.Struct(t)
}

// ValidateWithContext the struct.
func ValidateWithContext[T any](ctx context.Context, t *T) error {
	return validate.StructCtx(ctx, t)
}

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

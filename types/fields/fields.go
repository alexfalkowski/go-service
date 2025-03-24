package fields

import (
	"context"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Register the validator.
func Register(v *validator.Validate) {
	validate = v
}

// Validate for type t and a tag (required).
func Validate[T any](t *T, tag string) error {
	return validate.Var(t, tag)
}

// ValidateWithContext for type t and a tag (required).
func ValidateWithContext[T any](ctx context.Context, t *T, tag string) error {
	return validate.VarCtx(ctx, t, tag)
}

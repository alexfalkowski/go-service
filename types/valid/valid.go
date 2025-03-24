package valid

import (
	"context"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

// Struct for type t.
func Struct[T any](ctx context.Context, t *T) error {
	return validate.StructCtx(ctx, t)
}

// Field for type t and a tag (required).
func Field[T any](ctx context.Context, t *T, tag string) error {
	return validate.VarCtx(ctx, t, tag)
}

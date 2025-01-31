package valid

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

// Struct for type t.
func Struct[T any](t *T) error {
	return validate.Struct(t)
}

// Field for type t and a tag (required).
func Field[T any](t *T, tag string) error {
	return validate.Var(t, tag)
}

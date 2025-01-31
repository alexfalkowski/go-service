package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

// ValidateStruct for type t.
func ValidateStruct[T any](t *T) error {
	return validate.Struct(t)
}

// ValidateField for type t and a tag (required).
func ValidateField[T any](t *T, tag string) error {
	return validate.Var(t, tag)
}

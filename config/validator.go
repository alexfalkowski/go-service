package config

import "github.com/go-playground/validator/v10"

// NewValidator constructs a Validator backed by go-playground/validator.
//
// It enables required-struct validation (validator.WithRequiredStructEnabled), which causes validation
// tags like `required` to be applied to nested struct fields in a more strict/consistent way.
//
// This constructor is typically wired via `config.Module` and consumed by `NewConfig[T]` to validate
// decoded configuration before returning it to the caller.
func NewValidator() *Validator {
	return &Validator{validator.New(validator.WithRequiredStructEnabled())}
}

// Validator wraps a go-playground validator instance.
//
// It is used by `NewConfig[T]` to validate decoded configuration structs. You may use the embedded
// `*validator.Validate` directly to register custom validations or to validate values manually.
type Validator struct {
	*validator.Validate
}

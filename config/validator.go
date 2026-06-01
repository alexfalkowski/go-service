package config

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/go-playground/validator/v10"
)

// NewValidator constructs a Validator backed by go-playground/validator.
//
// It enables required-struct validation (validator.WithRequiredStructEnabled), which causes validation
// tags like `required` to be applied to nested struct fields in a more strict/consistent way.
//
// This constructor is typically wired via `config.Module` and consumed by `NewConfig[T]` to validate
// decoded configuration before returning it to the caller.
func NewValidator() *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())
	runtime.Must(validate.RegisterValidation("config_size", validateConfigSize))

	return &Validator{validate}
}

func validateConfigSize(fl validator.FieldLevel) bool {
	field := fl.Field()
	if !field.CanInt() {
		return false
	}

	size := bytes.Size(field.Int())
	return size >= 0 && size <= bytes.MaxConfigSize
}

// Validator wraps a go-playground validator instance.
//
// It is used by `NewConfig[T]` to validate decoded configuration structs. You may use the embedded
// `*validator.Validate` directly to register custom validations or to validate values manually.
type Validator struct {
	*validator.Validate
}

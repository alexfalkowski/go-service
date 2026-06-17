package config

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/go-playground/validator/v10"
)

// FieldLevel aliases the upstream validator field-level interface for custom config validation rules.
type FieldLevel = validator.FieldLevel

// NewValidator constructs a Validator backed by go-playground/validator.
//
// It enables required-struct validation ([validator.WithRequiredStructEnabled]), which causes validation
// tags like `required` to be applied to nested struct fields in a more strict/consistent way.
//
// It also registers repository-owned validation tags:
//   - `config_size`: accepts integer-like byte sizes between 0 and [bytes.MaxConfigSize].
//   - `duration_second_precision`: accepts positive durations that are exact multiples of one second.
//
// This constructor is typically wired via [Module] and consumed by `NewConfig[T]` to validate
// decoded configuration before returning it to the caller.
func NewValidator() *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())
	runtime.Must(validate.RegisterValidation("config_size", validateConfigSize))
	runtime.Must(validate.RegisterValidation("duration_second_precision", validateDurationSecondPrecision))

	return &Validator{validate}
}

// Validator wraps a go-playground validator instance.
//
// It is used by `NewConfig[T]` to validate decoded configuration structs. You may use the embedded
// `*validator.Validate` directly to register custom validations or to validate values manually.
type Validator struct {
	*validator.Validate
}

func validateConfigSize(fl FieldLevel) bool {
	field := fl.Field()
	size := bytes.Size(field.Int())
	return size >= 0 && size <= bytes.MaxConfigSize
}

func validateDurationSecondPrecision(fl FieldLevel) bool {
	field := fl.Field()
	duration := time.Duration(field.Int())
	return duration > 0 && duration%time.Second == 0
}

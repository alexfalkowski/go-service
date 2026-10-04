package validate

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/go-playground/validator/v10"
)

// FieldLevel aliases the upstream validator field-level interface for custom validation rules.
type FieldLevel = validator.FieldLevel

// Func aliases the upstream validator function signature for a custom validation rule.
type Func = validator.Func

// Validation names a go-playground/validator tag together with the function that implements it.
type Validation struct {
	// Func implements the rule for Tag.
	Func Func
	// Tag is the validation tag name, for example `config_size`.
	Tag string
}

// ValidationResult tags a Validation into the shared "validations" group consumed by NewValidator.
//
// Packages that own a validation rule return ValidationResult from their own constructor so
// NewValidator can register that rule without needing to know about it by name.
type ValidationResult struct {
	di.Out

	// Validation is collected into the "validations" group.
	Validation Validation `group:"validations"`
}

// ValidatorParams defines dependencies for NewValidator.
type ValidatorParams struct {
	di.In

	// Validations are the validation rules contributed by other packages via ValidationResult.
	Validations []Validation `group:"validations"`
}

// NewValidator constructs a Validator backed by go-playground/validator.
//
// It enables required-struct validation ([validator.WithRequiredStructEnabled]), which causes
// validation tags like `required` to be applied to nested struct fields in a more strict/consistent
// way, and registers every Validation contributed to the "validations" group.
func NewValidator(params ValidatorParams) *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled(), validator.WithTagNameFuncBlankOmit())
	for _, v := range params.Validations {
		runtime.Must(validate.RegisterValidation(v.Tag, v.Func))
	}

	return &Validator{validate}
}

// Validator wraps a go-playground validator instance.
//
// It is used by `config.NewConfig[T]` to validate decoded configuration structs. You may use the
// embedded `*validator.Validate` directly to validate values manually.
type Validator struct {
	*validator.Validate
}

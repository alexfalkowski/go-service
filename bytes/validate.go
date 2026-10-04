package bytes

import "github.com/alexfalkowski/go-service/v2/config/validate"

// newValidation declares the `config_size` validation rule owned by this package.
func newValidation() validate.ValidationResult {
	return validate.ValidationResult{Validation: validate.Validation{Tag: "config_size", Func: validateConfigSize}}
}

func validateConfigSize(fl validate.FieldLevel) bool {
	return ValidateConfigSize(Size(fl.Field().Int())) == nil
}

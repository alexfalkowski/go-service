package time

import "github.com/alexfalkowski/go-service/v2/config/validate"

// newValidation declares the `duration_second_precision` validation rule owned by this package.
func newValidation() validate.ValidationResult {
	return validate.ValidationResult{
		Validation: validate.Validation{Tag: "duration_second_precision", Func: validateSecondPrecision},
	}
}

func validateSecondPrecision(fl validate.FieldLevel) bool {
	return ValidateSecondPrecision(Duration(fl.Field().Int())) == nil
}

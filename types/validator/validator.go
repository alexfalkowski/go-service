package validator

import "github.com/go-playground/validator/v10"

// Validator is an alias for go-playground validator.
type Validator = validator.Validate

// NewValidator using go-playground validator.
func NewValidator() *Validator {
	return validator.New(validator.WithRequiredStructEnabled())
}

package validator

import "github.com/go-playground/validator/v10"

// NewValidator using go-playground validator.
func NewValidator() *validator.Validate {
	return validator.New(validator.WithRequiredStructEnabled())
}

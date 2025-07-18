package config

import "github.com/go-playground/validator/v10"

// NewValidator using go-playground validator.
func NewValidator() *Validator {
	return &Validator{validator.New(validator.WithRequiredStructEnabled())}
}

// Validator is a wrapper for go-playground validator.
type Validator struct {
	*validator.Validate
}

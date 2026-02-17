package config

import "github.com/go-playground/validator/v10"

// NewValidator constructs a Validator backed by go-playground/validator.
func NewValidator() *Validator {
	return &Validator{validator.New(validator.WithRequiredStructEnabled())}
}

// Validator wraps a go-playground validator instance.
type Validator struct {
	*validator.Validate
}

package cmd

import (
	"github.com/alexfalkowski/go-service/marshaller"
)

// InputConfig for cmd.
type InputConfig struct {
	*Config
}

// NewInputConfig for cmd.
func NewInputConfig(factory *marshaller.Factory) (*InputConfig, error) {
	c, err := NewConfig(inputFlag, factory)

	return &InputConfig{Config: c}, err
}

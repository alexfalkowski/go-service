package cmd

import (
	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/marshaller"
)

// InputFlag for cmd.
var InputFlag = String()

// InputConfig for cmd.
type InputConfig struct {
	*Config
}

// NewInputConfig for cmd.
func NewInputConfig(factory *marshaller.Factory) (*InputConfig, error) {
	c, err := NewConfig(*InputFlag, factory)

	return &InputConfig{Config: c}, errors.Prefix("new input", err)
}

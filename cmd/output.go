package cmd

import (
	"github.com/alexfalkowski/go-service/marshaller"
)

// OutputFlag for cmd.
var OutputFlag = String()

// OutputConfig for cmd.
type OutputConfig struct {
	*Config
}

// NewOutputConfig for cmd.
func NewOutputConfig(factory *marshaller.Factory) (*OutputConfig, error) {
	c, err := NewConfig(*OutputFlag, factory)
	if err != nil {
		return nil, err
	}

	return &OutputConfig{Config: c}, nil
}

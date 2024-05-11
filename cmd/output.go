package cmd

import (
	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/flags"
	"github.com/alexfalkowski/go-service/marshaller"
)

// OutputFlag for cmd.
var OutputFlag = flags.String()

// OutputConfig for cmd.
type OutputConfig struct {
	*Config
}

// NewOutputConfig for cmd.
func NewOutputConfig(factory *marshaller.Factory) (*OutputConfig, error) {
	c, err := NewConfig(*OutputFlag, factory)

	return &OutputConfig{Config: c}, errors.Prefix("new output", err)
}

// RegisterInput for cmd.
func (c *Command) RegisterOutput(env string) {
	value := "env:" + env

	flags.StringVar(c.root, OutputFlag, "output", "o", value, "input config location (format kind:location, default "+value+")")
}

package cmd

import (
	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/flags"
	"github.com/alexfalkowski/go-service/marshaller"
)

// InputFlag for cmd.
var InputFlag = flags.String()

// InputConfig for cmd.
type InputConfig struct {
	*Config
}

// NewInputConfig for cmd.
func NewInputConfig(factory *marshaller.Factory) (*InputConfig, error) {
	c, err := NewConfig(*InputFlag, factory)

	return &InputConfig{Config: c}, errors.Prefix("new input", err)
}

// RegisterInput for cmd.
func (c *Command) RegisterInput(env string) {
	value := "env:" + env

	flags.StringVar(c.root, InputFlag, "input", "i", value, "input config location (format kind:location, default "+value+")")
}

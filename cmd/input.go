package cmd

import (
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
func NewInputConfig(factory *marshaller.Map) *InputConfig {
	c := NewConfig(*InputFlag, factory)

	return &InputConfig{Config: c}
}

// RegisterInput for cmd.
func (c *Command) RegisterInput(value string) {
	flags.StringVar(c.root, InputFlag, "input", "i", value, "input config location (format kind:location, default "+value+")")
}

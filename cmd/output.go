package cmd

import (
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
func NewOutputConfig(factory *marshaller.Map) *OutputConfig {
	c := NewConfig(*OutputFlag, factory)

	return &OutputConfig{Config: c}
}

// RegisterInput for cmd.
func (c *Command) RegisterOutput(value string) {
	flags.StringVar(c.root, OutputFlag, "output", "o", value, "input config location (format kind:location, default "+value+")")
}

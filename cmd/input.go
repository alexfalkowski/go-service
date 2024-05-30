package cmd

import (
	"github.com/alexfalkowski/go-service/flags"
	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/spf13/cobra"
)

// InputFlag for cmd.
var InputFlag = flags.String()

// InputConfig for cmd.
type InputConfig struct {
	*Config
}

// NewInputConfig for cmd.
func NewInputConfig(mm *marshaller.Map) *InputConfig {
	c := NewConfig(*InputFlag, mm)

	return &InputConfig{Config: c}
}

// RegisterInput for cmd.
func (c *Command) RegisterInput(cmd *cobra.Command, value string) {
	flags.StringVar(cmd, InputFlag, "input", "i", value, "input config location (format kind:location)")
}

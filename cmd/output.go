package cmd

import (
	"github.com/alexfalkowski/go-service/flags"
	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/spf13/cobra"
)

// OutputFlag for cmd.
var OutputFlag = flags.String()

// OutputConfig for cmd.
type OutputConfig struct {
	*Config
}

// NewOutputConfig for cmd.
func NewOutputConfig(mm *marshaller.Map) *OutputConfig {
	c := NewConfig(*OutputFlag, mm)

	return &OutputConfig{Config: c}
}

// RegisterInput for cmd.
func (c *Command) RegisterOutput(cmd *cobra.Command, value string) {
	flags.StringVar(cmd, OutputFlag, "output", "o", value, "output config location (format kind:location)")
}

package cmd

import (
	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/flags"
	"github.com/alexfalkowski/go-service/os"
)

// InputFlag for cmd.
var InputFlag = flags.String()

// InputConfig for cmd.
type InputConfig struct {
	*Config
}

// NewInputConfig for cmd.
func NewInputConfig(enc *encoding.Map, fs os.FileSystem) *InputConfig {
	c := NewConfig(*InputFlag, enc, fs)

	return &InputConfig{Config: c}
}

// RegisterInput for cmd.
func (c *Command) RegisterInput(flags *flags.FlagSet, value string) {
	flags.StringVarP(InputFlag, "input", "i", value, "input config location (format kind:location)")
}

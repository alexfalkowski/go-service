package cmd

import (
	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/types/ptr"
	"github.com/leaanthony/clir"
)

// InputFlag for cmd.
var InputFlag *string

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
func (c *Command) RegisterInput(command *clir.Command, value string) {
	InputFlag = ptr.Value(value)
	command.StringFlag("input", "input config location (format kind:location)", InputFlag)
}

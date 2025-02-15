package cmd

import (
	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/types/ptr"
	"github.com/leaanthony/clir"
)

// OutputFlag for cmd.
var OutputFlag *string

// OutputConfig for cmd.
type OutputConfig struct {
	*Config
}

// NewOutputConfig for cmd.
func NewOutputConfig(enc *encoding.Map, fs os.FileSystem) *OutputConfig {
	c := NewConfig(*OutputFlag, enc, fs)

	return &OutputConfig{Config: c}
}

// RegisterInput for cmd.
func (c *Command) RegisterOutput(command *clir.Command, value string) {
	OutputFlag = ptr.Value(value)
	command.StringFlag("output", "output config location (format kind:location)", OutputFlag)
}

package cmd

import (
	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/flags"
	"github.com/alexfalkowski/go-service/os"
)

// InputConfig for cmd.
type InputConfig struct {
	*Config
}

// NewInputConfig for cmd.
func NewInputConfig(set *flags.FlagSet, enc *encoding.Map, fs os.FileSystem) *InputConfig {
	input, _ := set.GetString("input")
	config := NewConfig(input, enc, fs)

	return &InputConfig{Config: config}
}

package cmd

import (
	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/os"
)

// NewInputConfig for cmd.
func NewInputConfig(set *FlagSet, enc *encoding.Map, fs os.FileSystem) *InputConfig {
	input, _ := set.GetString("input")
	config := NewConfig(input, enc, fs)

	return &InputConfig{Config: config}
}

// InputConfig for cmd.
type InputConfig struct {
	*Config
}

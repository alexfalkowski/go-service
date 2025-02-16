package cmd

import (
	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/os"
)

// OutputConfig for cmd.
type OutputConfig struct {
	*Config
}

// NewOutputConfig for cmd.
func NewOutputConfig(set *FlagSet, enc *encoding.Map, fs os.FileSystem) *OutputConfig {
	output, _ := set.GetString("output")
	config := NewConfig(output, enc, fs)

	return &OutputConfig{Config: config}
}

package cmd

import (
	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
)

// NewOutputConfig for cmd.
func NewOutputConfig(name env.Name, set *FlagSet, enc *encoding.Map, fs *os.FS) *OutputConfig {
	return &OutputConfig{
		Config: NewConfig(name, set.GetOutput(), enc, fs),
	}
}

// OutputConfig for cmd.
type OutputConfig struct {
	*Config
}

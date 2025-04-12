package cmd

import (
	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
)

// NewInputConfig for cmd.
func NewInputConfig(name env.Name, set *FlagSet, enc *encoding.Map, fs os.FileSystem) *InputConfig {
	return &InputConfig{
		Config: NewConfig(name, set.GetInput(), enc, fs),
	}
}

// InputConfig for cmd.
type InputConfig struct {
	*Config
}

package cmd

import (
	"github.com/alexfalkowski/go-service/cmd/config"
	"github.com/alexfalkowski/go-service/cmd/flag"
	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
)

// NewInputConfig for cmd.
func NewInputConfig(name env.Name, set *flag.FlagSet, enc *encoding.Map, fs *os.FS) *InputConfig {
	return &InputConfig{config.NewConfig(name, set.GetInput(), enc, fs)}
}

// InputConfig for cmd.
type InputConfig struct {
	*config.Config
}

package cli

import (
	"github.com/alexfalkowski/go-service/v2/cli/config"
	"github.com/alexfalkowski/go-service/v2/cli/flag"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/os"
)

// NewInputConfig for cmd.
func NewInputConfig(name env.Name, flags *flag.FlagSet, enc *encoding.Map, fs *os.FS) *InputConfig {
	return &InputConfig{config.NewConfig(name, flags.GetInput(), enc, fs)}
}

// InputConfig for cmd.
type InputConfig struct {
	*config.Config
}

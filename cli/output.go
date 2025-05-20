package cli

import (
	"github.com/alexfalkowski/go-service/v2/cli/config"
	"github.com/alexfalkowski/go-service/v2/cli/flag"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/os"
)

// NewOutputConfig for cmd.
func NewOutputConfig(name env.Name, set *flag.FlagSet, enc *encoding.Map, fs *os.FS) *OutputConfig {
	return &OutputConfig{config.NewConfig(name, set.GetOutput(), enc, fs)}
}

// OutputConfig for cmd.
type OutputConfig struct {
	*config.Config
}

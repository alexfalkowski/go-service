package cli

import (
	"github.com/alexfalkowski/go-service/v2/cli/flag"
	"github.com/alexfalkowski/go-service/v2/config/io"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// ErrNoEncoder for cli.
var ErrNoEncoder = errors.New("config: no encoder")

// NewConfig for cli.
func NewConfig(name env.Name, flags *flag.FlagSet, enc *encoding.Map, fs *os.FS) *Config {
	kind, location := strings.CutColon(flags.GetInput())
	source := io.NewSource(name, kind, location, fs)
	encoder := enc.Get(source.Kind())

	return &Config{source: source, encoder: encoder}
}

// Config for cli.
type Config struct {
	encoder encoding.Encoder
	source  io.Source
}

// Decode for config.
func (c *Config) Decode(v any) error {
	if c.encoder == nil {
		return ErrNoEncoder
	}

	reader := c.source.Reader()
	defer reader.Close()

	return c.prefix(c.encoder.Decode(reader, v))
}

func (c *Config) prefix(err error) error {
	return errors.Prefix("config", err)
}

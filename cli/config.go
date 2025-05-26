package cli

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
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
	reader := io.NewReader(name, kind, location, fs)
	encoder := enc.Get(reader.Kind())

	return &Config{reader: reader, encoder: encoder}
}

// Config for cli.
type Config struct {
	encoder encoding.Encoder
	reader  io.Reader
}

// Decode for config.
func (c *Config) Decode(v any) error {
	data, err := c.reader.Read()
	if err != nil {
		return c.prefix(err)
	}

	if c.encoder == nil {
		return ErrNoEncoder
	}

	return c.prefix(c.encoder.Decode(bytes.NewBuffer(data), v))
}

func (c *Config) prefix(err error) error {
	return errors.Prefix("config", err)
}

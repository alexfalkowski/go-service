package config

import (
	"errors"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/config/io"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/env"
	se "github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// ErrNoEncoder for cmd.
var ErrNoEncoder = errors.New("config: no encoder")

// NewConfig for cmd.
func NewConfig(name env.Name, arg string, enc *encoding.Map, fs *os.FS) *Config {
	kind, location := strings.CutColon(arg)
	reader := io.NewReader(name, kind, location, fs)
	encoder := enc.Get(reader.Kind())

	return &Config{reader: reader, encoder: encoder}
}

// Config for cmd.
type Config struct {
	encoder encoding.Encoder
	reader  io.Reader
}

// Decode for config.
func (c *Config) Decode(data any) error {
	bts, err := c.reader.Read()
	if err != nil {
		return c.prefix(err)
	}

	if c.encoder == nil {
		return ErrNoEncoder
	}

	return c.prefix(c.encoder.Decode(bytes.NewBuffer(bts), data))
}

func (c *Config) prefix(err error) error {
	return se.Prefix("config", err)
}

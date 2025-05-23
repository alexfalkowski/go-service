package config

import (
	"errors"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/cli/flag"
	"github.com/alexfalkowski/go-service/v2/config/io"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/env"
	se "github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
)

// ErrNoEncoder for cmd.
var ErrNoEncoder = errors.New("config: no encoder")

// NewConfig for cmd.
func NewConfig(name env.Name, arg string, enc *encoding.Map, fs *os.FS) *Config {
	kind, location := flag.SplitFlag(arg)
	rw := io.NewReadWriter(name, kind, location, fs)
	encoder := enc.Get(rw.Kind())

	return &Config{rw: rw, encoder: encoder}
}

// Config for cmd.
type Config struct {
	encoder encoding.Encoder
	rw      io.ReaderWriter
}

// Kind of config.
func (c *Config) Kind() string {
	return c.rw.Kind()
}

// Decode for config.
func (c *Config) Decode(data any) error {
	bts, err := c.rw.Read()
	if err != nil {
		return c.prefix(err)
	}

	if c.encoder == nil {
		return ErrNoEncoder
	}

	return c.prefix(c.encoder.Decode(bytes.NewBuffer(bts), data))
}

// Write for config.
func (c *Config) Write(data []byte, mode os.FileMode) error {
	return c.prefix(c.rw.Write(data, mode))
}

func (c *Config) prefix(err error) error {
	return se.Prefix("config", err)
}

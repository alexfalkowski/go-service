package cmd

import (
	"bytes"
	"errors"

	"github.com/alexfalkowski/go-service/encoding"
	se "github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/os"
)

// ErrNoEncoder for cmd.
var ErrNoEncoder = errors.New("config: no encoder")

// Config for cmd.
type Config struct {
	enc encoding.Encoder
	rw  ReaderWriter
}

// NewConfig for cmd.
func NewConfig(flag string, enc *encoding.Map, fs os.FileSystem) *Config {
	kind, location := SplitFlag(flag)
	rw := NewReadWriter(kind, location, fs)
	m := enc.Get(rw.Kind())

	return &Config{rw: rw, enc: m}
}

// Kind of config.
func (c *Config) Kind() string {
	return c.rw.Kind()
}

// Decode for config.
func (c *Config) Decode(data any) error {
	bts, err := c.rw.Read()
	if err != nil {
		return se.Prefix("config", err)
	}

	if c.enc == nil {
		return ErrNoEncoder
	}

	return se.Prefix("config", c.enc.Decode(bytes.NewBufferString(bts), data))
}

// Write for config.
func (c *Config) Write(data string, mode os.FileMode) error {
	return se.Prefix("config", c.rw.Write(data, mode))
}

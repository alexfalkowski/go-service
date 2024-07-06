package cmd

import (
	"io/fs"

	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/errors"
)

// Config for cmd.
type Config struct {
	m  encoding.Marshaller
	rw ReaderWriter
}

// NewConfig for cmd.
func NewConfig(flag string, enc *encoding.Map) *Config {
	k, l := SplitFlag(flag)
	rw := NewReadWriter(k, l)
	m := enc.Get(rw.Kind())

	return &Config{rw: rw, m: m}
}

// Kind of config.
func (c *Config) Kind() string {
	return c.rw.Kind()
}

// Unmarshal for config.
func (c *Config) Unmarshal(data any) error {
	d, err := c.rw.Read()
	if err != nil {
		return errors.Prefix("unmarshal config", err)
	}

	return errors.Prefix("unmarshal config", c.m.Unmarshal(d, data))
}

// Write for config.
func (c *Config) Write(data []byte, mode fs.FileMode) error {
	return errors.Prefix("write config", c.rw.Write(data, mode))
}

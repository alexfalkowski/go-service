package cmd

import (
	"io/fs"

	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/marshaller"
)

// ReaderWriter for cmd.
type ReaderWriter interface {
	// Read bytes.
	Read() ([]byte, error)

	// Write bytes with files's mode.
	Write(data []byte, mode fs.FileMode) error

	// Kind of read writer.
	Kind() string
}

// Config for cmd.
type Config struct {
	m  marshaller.Marshaller
	rw ReaderWriter
}

// NewConfig for cmd.
func NewConfig(flag string, factory *marshaller.Factory) (*Config, error) {
	k, l := SplitFlag(flag)
	rw := NewReadWriter(k, l)
	m, err := factory.Create(rw.Kind())

	return &Config{rw: rw, m: m}, errors.Prefix("new config", err)
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

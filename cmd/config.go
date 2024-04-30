package cmd

import (
	"io/fs"
	"strings"

	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/marshaller"
)

var inputFlag string

// ReaderWriter for cmd.
type ReaderWriter interface {
	Read() ([]byte, error)
	Write(data []byte, mode fs.FileMode) error
	Kind() string
}

// Config for cmd.
type Config struct {
	m  marshaller.Marshaller
	rw ReaderWriter
}

// NewConfig for cmd.
func NewConfig(flag string, factory *marshaller.Factory) (*Config, error) {
	k, l := splitFlag(flag)
	rw := readWriter(k, l)
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

func splitFlag(f string) (string, string) {
	c := strings.Split(f, ":")

	if len(c) != 2 {
		return "env", "CONFIG_FILE"
	}

	return c[0], c[1]
}

func readWriter(k, l string) ReaderWriter {
	if k == "file" {
		return NewFile(l)
	}

	return NewENV(l)
}

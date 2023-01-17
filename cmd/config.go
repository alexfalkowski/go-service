package cmd

import (
	"errors"
	"io/fs"
	"strings"

	"github.com/alexfalkowski/go-service/marshaller"
)

var (
	// ErrInvalidKind for cmd.
	ErrInvalidKind = errors.New("invalid kind")
	inputFlag      string
)

// ReaderWriter for cmd.
type ReaderWriter interface {
	Read() ([]byte, error)
	Write([]byte, fs.FileMode) error
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

	var rw ReaderWriter

	switch k {
	case "env":
		rw = NewENV(l)
	case "file":
		rw = NewFile(l)
	}

	if rw == nil {
		return nil, ErrInvalidKind
	}

	m, err := factory.Create(rw.Kind())
	if err != nil {
		return nil, err
	}

	return &Config{rw: rw, m: m}, nil
}

// Unmarshal for config.
func (c *Config) Unmarshal(data any) error {
	d, err := c.rw.Read()
	if err != nil {
		return err
	}

	return c.m.Unmarshal(d, data)
}

// Write for config.
func (c *Config) Write(data []byte, mode fs.FileMode) error {
	return c.rw.Write(data, mode)
}

func splitFlag(f string) (string, string) {
	c := strings.Split(f, ":")

	if len(c) != 2 {
		return "env", "CONFIG_FILE"
	}

	return c[0], c[1]
}

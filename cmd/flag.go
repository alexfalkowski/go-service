package cmd

import (
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexfalkowski/go-service/file"
	"github.com/alexfalkowski/go-service/marshaller"
)

var (
	// ErrInvalidKind for cmd.
	ErrInvalidKind = errors.New("invalid kind")

	// ConfigFlag for cmd.
	ConfigFlag string
)

// Config for cmd.
type Config struct {
	kind    string
	Data    []byte
	factory *marshaller.Factory
}

// NewConfig for cmd.
func NewConfig(factory *marshaller.Factory) (*Config, error) {
	k, l := splitFlag()
	switch k {
	case "env":
		k := file.Extension(os.Getenv(l))
		d, err := os.ReadFile(filepath.Clean(os.Getenv(l)))

		return &Config{kind: k, Data: d, factory: factory}, err
	case "file":
		k := file.Extension(l)
		d, err := os.ReadFile(filepath.Clean(l))

		return &Config{kind: k, Data: d, factory: factory}, err
	case "mem":
		k, l = splitMemory(l)
		d, err := base64.StdEncoding.DecodeString(l)

		return &Config{kind: k, Data: d, factory: factory}, err
	}

	return nil, ErrInvalidKind
}

// Marshaller for cmd.
func (c *Config) Marshaller() (marshaller.Marshaller, error) {
	return c.factory.Create(c.kind)
}

func splitFlag() (string, string) {
	c := strings.Split(ConfigFlag, ":")

	if len(c) != 2 {
		return "env", "CONFIG_FILE"
	}

	return c[0], c[1]
}

func splitMemory(l string) (string, string) {
	c := strings.Split(l, "=>")

	if len(c) != 2 {
		return "yaml", l
	}

	return c[0], c[1]
}

package cmd

import (
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ErrInvalidKind for cmd.
	ErrInvalidKind = errors.New("invalid kind")

	// ConfigFlag for cmd.
	ConfigFlag string
)

// Config for cmd.
type Config struct {
	Data []byte
}

// NewConfig for cmd.
func NewConfig() (*Config, error) {
	k, l := split()
	switch k {
	case "env":
		d, err := os.ReadFile(filepath.Clean(os.Getenv(l)))

		return &Config{Data: d}, err
	case "file":
		d, err := os.ReadFile(filepath.Clean(l))

		return &Config{Data: d}, err
	case "mem":
		d, err := base64.StdEncoding.DecodeString(l)

		return &Config{Data: d}, err
	}

	return nil, ErrInvalidKind
}

func split() (string, string) {
	c := strings.Split(ConfigFlag, ":")

	if len(c) != 2 {
		return "env", "CONFIG_FILE"
	}

	return c[0], c[1]
}

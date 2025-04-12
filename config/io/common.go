package io

import (
	"path/filepath"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/strings"
)

// NewCommon for io.
func NewCommon(name env.Name, fs os.FileSystem, rw ReaderWriter) *Common {
	return &Common{name: name, fs: fs, rw: rw}
}

// Common for io.
type Common struct {
	fs   os.FileSystem
	rw   ReaderWriter
	name env.Name
}

// It will first try to read from the underlying read writer, otherwise it will read a configuration file in commonly defined locations.
func (c *Common) Read() ([]byte, error) {
	bytes, err := c.rw.Read()
	if err != nil {
		name := c.name.String()
		file := name + ".yml"
		dirs := []string{
			os.ExecutableDir(),
			filepath.Join(os.UserHomeDir(), ".config", name),
			"/etc/" + name,
		}

		for _, dir := range dirs {
			name := filepath.Join(dir, file)
			if c.fs.PathExists(name) {
				return c.fs.ReadFile(name)
			}
		}

		return nil, err
	}

	return bytes, nil
}

// Write to the underlying read writer.
func (c *Common) Write(data []byte, mode os.FileMode) error {
	return c.rw.Write(data, mode)
}

// Kind from the underlying read writer, otherwise YAML.
func (c *Common) Kind() string {
	kind := c.rw.Kind()
	if strings.IsEmpty(kind) {
		return "yml"
	}

	return kind
}

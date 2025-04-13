package io

import (
	"path/filepath"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/strings"
)

// NewCommon for io.
func NewCommon(name env.Name, fs os.FileSystem, rw ReaderWriter) *Common {
	location, kind := common(name)

	return &Common{location: location, kind: kind, fs: fs, rw: rw}
}

// Common for io.
type Common struct {
	fs       os.FileSystem
	rw       ReaderWriter
	location string
	kind     string
}

// Valid checks if the location is present.
func (c *Common) Valid() bool {
	if c.rw != nil {
		return c.rw.Valid()
	}

	return c.fs.PathExists(c.location)
}

// It will first try to read from the underlying read writer, otherwise it will read a configuration file in commonly defined locations.
func (c *Common) Read() ([]byte, error) {
	if c.valid() {
		return c.rw.Read()
	}

	if strings.IsEmpty(c.location) {
		return nil, ErrLocationMissing
	}

	return c.fs.ReadFile(c.location)
}

// Write to the underlying read writer.
func (c *Common) Write(data []byte, mode os.FileMode) error {
	if c.valid() {
		return c.rw.Write(data, mode)
	}

	if strings.IsEmpty(c.location) {
		return ErrLocationMissing
	}

	return c.fs.WriteFile(c.location, data, mode)
}

// Kind from the underlying read writer, otherwise YAML.
func (c *Common) Kind() string {
	if c.valid() {
		return c.rw.Kind()
	}

	if strings.IsEmpty(c.kind) {
		return "yml"
	}

	return c.kind
}

func (c *Common) valid() bool {
	return c.rw != nil && c.rw.Valid()
}

func common(name env.Name) (string, string) {
	extensions := []string{".yml", ".toml", ".json"}
	for _, extension := range extensions {
		n := name.String()
		file := n + extension
		dirs := []string{
			os.ExecutableDir(),
			filepath.Join(os.UserHomeDir(), ".config", n),
			"/etc/" + n,
		}

		for _, dir := range dirs {
			name := filepath.Join(dir, file)
			if os.PathExists(name) {
				return name, os.PathExtension(name)
			}
		}
	}

	return "", ""
}

package cmd

import (
	"encoding/base64"
	"io/fs"
	"strings"
)

// MEM for cmd.
type MEM struct {
	location string
}

// NewMEM for cmd.
func NewMEM(location string) *MEM {
	return &MEM{location: location}
}

// Read for MEM.
func (m *MEM) Read() ([]byte, error) {
	_, l := m.split()

	return base64.StdEncoding.DecodeString(l)
}

// Write for MEM.
func (m *MEM) Write(data []byte, mode fs.FileMode) error {
	return nil
}

// Write for MEM.
func (m *MEM) Kind() string {
	k, _ := m.split()

	return k
}

func (m *MEM) split() (string, string) {
	c := strings.Split(m.location, "=>")

	if len(c) != 2 {
		return "yaml", m.location
	}

	return c[0], c[1]
}

package cmd

import (
	"io/fs"
)

// None for cmd.
type None struct{}

// NewNone for cmd.
func NewNone() *None {
	return &None{}
}

// Read for none.
func (*None) Read() ([]byte, error) {
	return nil, nil
}

// Write for none.
func (*None) Write(_ []byte, _ fs.FileMode) error {
	return nil
}

// Kind for none.
func (*None) Kind() string {
	return "none"
}
